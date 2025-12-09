package services

import (
	"context"
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"github.com/vviveksharma/auth/db"
	"github.com/vviveksharma/auth/internal/models"
	"github.com/vviveksharma/auth/internal/repo"
	dbmodels "github.com/vviveksharma/auth/models"
	"github.com/vviveksharma/auth/queue"
)

type MessageService interface {
	CreateMessage(req *models.CreateMessageRequest, ctx context.Context) (resp *models.CreateMessageResponse, err error)
	GetStatus(messageId uuid.UUID, ctx context.Context) (resp *models.GetMessageStatusResponse, err error)
	ListMessages(email string, ctx context.Context) (resp []*models.ListMessageStatusResponse, err error)
}

type Message struct {
	MessageRepo  repo.MessageRepositoryInterface
	UserRepo     repo.UserRepositoryInterface
	QConn        *amqp.Connection
	Queue        amqp.Queue
	QueueService queue.IQueueService
}

func NewMessageService(qu amqp.Queue, conn *amqp.Connection) (MessageService, error) {
	ser := &Message{}
	err := ser.SetupRepo()
	ser.QConn = conn
	ser.Queue = qu
	if err != nil {
		return nil, err
	}
	return ser, nil
}

func (m *Message) SetupRepo() error {
	var err error
	message, err := repo.NewMessageRepository(db.DB)
	if err != nil {
		return err
	}
	m.MessageRepo = message
	user, err := repo.NewUserRepository(db.DB)
	if err != nil {
		return err
	}
	m.UserRepo = user
	queue, err := queue.NewQueueRequest()
	if err != nil {
		return err
	}
	m.QueueService = queue
	return nil
}

func (m *Message) CreateMessage(req *models.CreateMessageRequest, ctx context.Context) (resp *models.CreateMessageResponse, err error) {
	log.Printf("CreateMessage: Starting function execution")

	// Check context values
	tenant_id_val := ctx.Value("tenant_id")
	user_id_val := ctx.Value("user_id")

	tenant_id := tenant_id_val.(string)
	user_id := user_id_val.(string)

	log.Printf("CreateMessage: tenant_id=%s, user_id=%s", tenant_id, user_id)

	// Check if request is nil
	if req == nil {
		log.Printf("CreateMessage: request is nil")
		return nil, &dbmodels.ServiceResponse{
			Code:    422,
			Message: "request cannot be nil",
		}
	}

	userResp, err := m.UserRepo.GetUserDetails(dbmodels.DBUser{
		Id:       uuid.MustParse(user_id),
		TenantId: uuid.MustParse(tenant_id),
	})
	if err != nil {
		log.Printf("CreateMessage: Error getting user details: %v", err)
		if err.Error() != "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "error while fetching the user details: " + err.Error(),
			}
		}
	}

	if userResp == nil {
		log.Printf("CreateMessage: userResp is nil")
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "user details not found",
		}
	}

	log.Printf("CreateMessage: User roles: %v", userResp.Roles)

	if slices.Contains(userResp.Roles, req.RequestedRole) {
		log.Printf("CreateMessage: User already has the requested role: %s", req.RequestedRole)
		return nil, &dbmodels.ServiceResponse{
			Code:    409,
			Message: fmt.Sprintf("Access denied: You already possess the '%s' role. Role escalation requests are only permitted for roles you do not currently hold. Please verify your existing permissions or contact your administrator if you believe this is an error.", req.RequestedRole),
		}
	}

	log.Printf("CreateMessage: User doesn't have the requested role, proceeding...")

	// Check if MessageRepo is nil
	if m.MessageRepo == nil {
		log.Printf("CreateMessage: MessageRepo is nil")
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "MessageRepo is not initialized",
		}
	}

	message, err := m.MessageRepo.GetMessageByConditions(dbmodels.DBMessage{
		UserEmail:     req.Email,
		RequestedRole: req.RequestedRole,
		TenantId:      uuid.MustParse(tenant_id),
	})
	if err != nil {
		log.Printf("CreateMessage: Error getting message by conditions: %v", err)
		if err.Error() != "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "error while getting the message requested inforation: " + err.Error(),
			}
		}
	}

	if message != nil && message.RequestedRole == req.RequestedRole && req.Email == message.UserEmail {
		log.Printf("CreateMessage: Duplicate request found - Role: %s, Email: %s, Status: %s", message.RequestedRole, message.UserEmail, message.Status)
		return nil, &dbmodels.ServiceResponse{
			Code:    409,
			Message: fmt.Sprintf("the role you requested is already has been queued and the correct status for you request is %s", message.Status),
		}
	}

	log.Printf("CreateMessage: Creating new message entry sending the message to the queue")

	err = m.QueueService.PublishMessageTask(m.Queue, m.QConn, dbmodels.DBMessage{
		UserEmail:     req.Email,
		TenantId:      uuid.MustParse(tenant_id),
		RequestedRole: req.RequestedRole,
		Status:        "pending",
		RequestAt:     time.Now(),
	})

	if err != nil {
		log.Printf("CreateMessage: Error creating message: %v", err)
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while creating a new entry in the message db: " + err.Error(),
		}
	}

	log.Printf("CreateMessage: Message created successfully")
	return &models.CreateMessageResponse{
		Message: "message sent to the tenant-admin",
	}, nil
}

func (m *Message) ListMessages(email string, ctx context.Context) (resp []*models.ListMessageStatusResponse, err error) {
	tenant_id := ctx.Value("tenant_id").(string)
	message, err := m.MessageRepo.GetMessagesForUser(dbmodels.DBMessage{
		UserEmail: email,
		TenantId:  uuid.MustParse(tenant_id),
	})
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "Messages not found. Please verify the message ID or submit a new request.",
			}
		}
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while getting the message requested inforation:" + err.Error(),
		}
	}
	for _, req := range message {
		message := models.ListMessageStatusResponse{
			MessageId:     req.Id.String(),
			Status:        req.Status,
			RequestedRole: req.RequestedRole,
		}
		resp = append(resp, &message)
	}
	fmt.Println("the response: ", resp)
	return resp, nil
}

func (m *Message) GetStatus(messageId uuid.UUID, ctx context.Context) (resp *models.GetMessageStatusResponse, err error) {
	tenant_id := ctx.Value("tenant_id").(string)
	message, err := m.MessageRepo.GetMessageByConditions(dbmodels.DBMessage{
		Id:       messageId,
		TenantId: uuid.MustParse(tenant_id),
	})
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "Message not found. Please verify the message ID or submit a new request.",
			}
		}
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while getting the message requested inforation:" + err.Error(),
		}
	}
	response, err := m.MessageRepo.GetStatus(messageId, message.TenantId)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while getting the message status",
		}
	}
	return &models.GetMessageStatusResponse{
		Status: "the status for your request is : " + response,
	}, nil
}
