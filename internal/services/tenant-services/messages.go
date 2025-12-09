package tenantservices

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/db"
	"github.com/vviveksharma/auth/internal/models"
	tenantmodels "github.com/vviveksharma/auth/internal/models/tenantModels"
	"github.com/vviveksharma/auth/internal/pagination"
	"github.com/vviveksharma/auth/internal/repo"
	tenantrepo "github.com/vviveksharma/auth/internal/repo/tenantRepo"
	dbmodels "github.com/vviveksharma/auth/models"
)

type ITenantMessageService interface {
	ListMessageRequest(ctx context.Context, page int, page_size int, status string) (resp *pagination.PaginatedResponse[*tenantmodels.TenantMessageResponseBoy], err error)
	ApproveRequest(ctx context.Context, messageId uuid.UUID) (resp *tenantmodels.TenantApproveMessageResponseBody, err error)
	RejectRequest(ctx context.Context, messageId uuid.UUID) (resp *tenantmodels.TenantRejectMessageResponseBody, err error)
}

type TenanatMessageService struct {
	MessageRepo       repo.MessageRepositoryInterface
	TenantMessageRepo tenantrepo.TenantMessageRepositoryInterface
	UserRepo          repo.UserRepositoryInterface
	RoleRepo          repo.RoleRepositoryInterface
}

func NewTenantMessageService() (ITenantMessageService, error) {
	ser := &TenanatMessageService{}
	err := ser.SetupRepo()
	if err != nil {
		return nil, err
	}
	return ser, nil
}

func (tm *TenanatMessageService) SetupRepo() error {
	var err error
	message, err := repo.NewMessageRepository(db.DB)
	if err != nil {
		return err
	}
	tm.MessageRepo = message
	tmessage, err := tenantrepo.NewTenantMessageRepository(db.DB)
	if err != nil {
		return err
	}
	tm.TenantMessageRepo = tmessage
	user, err := repo.NewUserRepository(db.DB)
	if err != nil {
		return err
	}
	tm.UserRepo = user
	role, err := repo.NewRoleRepository(db.DB)
	if err != nil {
		return err
	}
	tm.RoleRepo = role
	return nil
}

func (tm *TenanatMessageService) ListMessageRequest(ctx context.Context, page int, page_size int, status string) (resp *pagination.PaginatedResponse[*tenantmodels.TenantMessageResponseBoy], err error) {
	tenantId := ctx.Value("tenant_id").(string)
	log.Println("the tenant_id: ", tenantId)
	messages, totalCount, err := tm.TenantMessageRepo.ListMessages(uuid.MustParse(tenantId), page, page_size, status)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while rendering the tenant messages based filtered response: " + err.Error(),
		}
	}
	var response []*tenantmodels.TenantMessageResponseBoy
	for _, message := range messages {
		formattedTime := message.RequestAt.Format("Jan 2, 2006")
		tempMessage := &tenantmodels.TenantMessageResponseBoy{
			MessageId:     message.Id,
			UserEmail:     message.UserEmail,
			CurrentRole:   message.CurrentRole,
			RequestedRole: message.RequestedRole,
			Status:        message.Status,
			RequestAt:     formattedTime,
		}
		response = append(response, tempMessage)
	}
	totalPages := int64(0)
	if page_size > 0 {
		totalPages = (totalCount + int64(page_size) - 1) / int64(page_size)
	}
	hasNext := int64(page) < totalPages
	hasPrev := page > 1
	paginatedResponse := &pagination.PaginatedResponse[*tenantmodels.TenantMessageResponseBoy]{
		Data: response,
		Pagination: pagination.PaginationMeta{
			Page:       page,
			PageSize:   page_size,
			TotalPages: int(totalPages),
			TotalItems: int(totalCount),
			HasNext:    hasNext,
			HasPrev:    hasPrev,
		},
	}

	return paginatedResponse, nil
}

func (tm *TenanatMessageService) ApproveRequest(ctx context.Context, messageId uuid.UUID) (resp *tenantmodels.TenantApproveMessageResponseBody, err error) {
	tenantIdVal := ctx.Value("tenant_id")
	if tenantIdVal == nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    401,
			Message: "tenant_id not found in context. Authentication required.",
		}
	}
	tenantId, ok := tenantIdVal.(string)
	if !ok {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "invalid tenant_id format in context",
		}
	}

	status, err := tm.MessageRepo.GetMessageByConditions(dbmodels.DBMessage{
		TenantId: uuid.MustParse(tenantId),
		Id:       messageId,
		Action:   false,
	})
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "Message not found or already processed",
			}
		}
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while getting the status of the given messageId: " + err.Error(),
		}
	}
	if status.Status == "approved" {
		return nil, &dbmodels.ServiceResponse{
			Code:    409,
			Message: "The request is already approved",
		}
	}
	if status.Status == "rejected" {
		return nil, &dbmodels.ServiceResponse{
			Code:    409,
			Message: "Cannot approve a rejected request",
		}
	}
	err = tm.TenantMessageRepo.ApproveMessage(uuid.MustParse(tenantId), messageId)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while updating the status of the request",
		}
	}

	// After apporving lets grant the user that role as well
	userDetails, err := tm.UserRepo.GetUserDetails(dbmodels.DBUser{
		TenantId: uuid.MustParse(tenantId),
		Email:    status.UserEmail,
	})
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    400,
				Message: "user with this email doesn't exists",
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("error while fetching the userdetails with the email %s", status.UserEmail),
			}
		}
	}
	// Fetching the roleid if the system role then we will provide that role else will check on the tenant basis
	roleId, err := tm.RoleRepo.FindRoleId(status.RequestedRole)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    400,
				Message: "role that is request can't be granted as it doesn't exists",
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("error while finding the roleId for this given role %s", err.Error()),
			}
		}
	}
	if models.IsSystemRole(roleId) {
		fmt.Println("its the system role request proceeding with updating the role for the user")
		err := tm.UserRepo.UpdateUserRoles(userDetails.Id, status.RequestedRole)
		if err != nil {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "erorr while assigning the role to the user",
			}
		}
	} else {
		roleDetails, err := tm.RoleRepo.GetRolesDetails(&dbmodels.DBRoles{
			RoleId:   roleId,
			TenantId: uuid.MustParse(tenantId),
		})
		if err != nil {
			if err.Error() == "record not found" {
				return nil, &dbmodels.ServiceResponse{
					Code:    400,
					Message: "unable to find the given role in the db",
				}
			} else {
				return nil, &dbmodels.ServiceResponse{
					Code:    500,
					Message: "error while getting the roledetails" + err.Error(),
				}
			}
		}
		fmt.Println("the role details: ", roleDetails)
		err = tm.UserRepo.UpdateUserRoles(userDetails.Id, status.RequestedRole)
		if err != nil {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "erorr while assigning the role to the user",
			}

		}
	}
	return &tenantmodels.TenantApproveMessageResponseBody{
		Message: "Apprvoed the request for the messageId " + messageId.String(),
	}, nil
}

func (tm *TenanatMessageService) RejectRequest(ctx context.Context, messageId uuid.UUID) (resp *tenantmodels.TenantRejectMessageResponseBody, err error) {
	tenantIdVal := ctx.Value("tenant_id")
	if tenantIdVal == nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    401,
			Message: "tenant_id not found in context. Authentication required.",
		}
	}
	tenantId, ok := tenantIdVal.(string)
	if !ok {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "invalid tenant_id format in context",
		}
	}

	status, err := tm.MessageRepo.GetMessageByConditions(dbmodels.DBMessage{
		TenantId: uuid.MustParse(tenantId),
		Id:       messageId,
		Action:   false,
	})
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "Message not found or already processed",
			}
		}
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while getting the status of the given messageId: " + err.Error(),
		}
	}
	if status.Status == "rejected" {
		return nil, &dbmodels.ServiceResponse{
			Code:    409,
			Message: "The request is already rejected",
		}
	}
	if status.Status == "approved" {
		return nil, &dbmodels.ServiceResponse{
			Code:    409,
			Message: "Cannot reject an approved request",
		}
	}
	err = tm.TenantMessageRepo.RejectMessage(uuid.MustParse(tenantId), messageId)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while updating the status of the request",
		}
	}
	return &tenantmodels.TenantRejectMessageResponseBody{
		Message: "Rejected the request for the messageId " + messageId.String(),
	}, nil
}
