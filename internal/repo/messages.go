package repo

import (
	"github.com/google/uuid"
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

type MessageRepositoryInterface interface {
	Create(req *models.DBMessage) error
	GetStatus(messageId uuid.UUID, tenantId uuid.UUID) (string, error)
	GetMessageByConditions(conditions models.DBMessage) (resp *models.DBMessage, err error)
	GetMessagesForUser(conditions models.DBMessage) (resp []*models.DBMessage, err error)
}

type MessageRepository struct {
	DB *gorm.DB
}

func NewMessageRepository(db *gorm.DB) (MessageRepositoryInterface, error) {
	return &MessageRepository{DB: db}, nil
}

func (m *MessageRepository) Create(req *models.DBMessage) error {
	transaction := m.DB.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	create := transaction.Create(&req)
	if create.Error != nil {
		return create.Error
	}
	transaction.Commit()
	return nil
}

func (m *MessageRepository) GetStatus(messageId uuid.UUID, tenantId uuid.UUID) (string, error) {
	var message models.DBMessage
	err := m.DB.Where("id = ? AND tenant_id = ?", messageId, tenantId).First(&message).Error
	if err != nil {
		return "", err
	}
	return message.Status, nil
}

func (m *MessageRepository) GetMessageByConditions(conditions models.DBMessage) (resp *models.DBMessage, err error) {
	transaction := m.DB.Begin()
	if transaction.Error != nil {
		return nil, transaction.Error
	}
	defer transaction.Rollback()
	var message models.DBMessage
	getErr := transaction.Where(&conditions).First(&message)
	if getErr.Error != nil {
		return nil, getErr.Error
	}
	return &message, nil
}

func (m *MessageRepository) GetMessagesForUser(conditions models.DBMessage) (resp []*models.DBMessage, err error) {
	transaction := m.DB.Begin()
	if transaction.Error != nil {
		return nil, transaction.Error
	}
	defer transaction.Rollback()
	var message []*models.DBMessage
	getErr := transaction.Where(&conditions).Find(&message)
	if getErr.Error != nil {
		return nil, getErr.Error
	}
	return message, nil
}
