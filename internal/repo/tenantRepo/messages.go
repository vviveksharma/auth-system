package tenantrepo

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

type TenantMessageRepositoryInterface interface {
	ListMessages(tenantId uuid.UUID, page int, pageSize int, status string) ([]*models.DBMessage, int64, error)
	ApproveMessage(tenantId uuid.UUID, messageId uuid.UUID) error
	RejectMessage(tenantId uuid.UUID, messageId uuid.UUID) error
}

type TenantMessageRepository struct {
	DB *gorm.DB
}

func NewTenantMessageRepository(db *gorm.DB) (TenantMessageRepositoryInterface, error) {
	return &TenantMessageRepository{DB: db}, nil
}

func (tm *TenantMessageRepository) ListMessages(tenantId uuid.UUID, page int, pageSize int, status string) ([]*models.DBMessage, int64, error) {
	var totalCount int64
	var messages []*models.DBMessage

	// Start transaction for data consistency
	transaction := tm.DB.Begin()
	if transaction.Error != nil {
		return nil, 0, transaction.Error
	}
	defer transaction.Rollback()

	// Base query - filter by tenant and exclude actioned/deleted messages
	baseQuery := transaction.Model(&models.DBMessage{}).
		Where("tenant_id = ?", tenantId)

	// Apply status filter if provided
	// Supports: "pending", "approved", "rejected", or empty for all
	if status != "" && (status == "pending" || status == "approved" || status == "rejected") {
		baseQuery = baseQuery.Where("status = ?", status)
	}

	// Count total matching records (for pagination metadata)
	if err := baseQuery.Count(&totalCount).Error; err != nil {
		log.Printf("Error counting tenant messages: %v", err)
		return nil, 0, fmt.Errorf("error counting messages for tenant: %v", err)
	}

	// If no messages found, return early
	if totalCount == 0 {
		if err := transaction.Commit().Error; err != nil {
			return nil, 0, err
		}
		return []*models.DBMessage{}, 0, nil
	}

	// Calculate offset for pagination
	offset := (page - 1) * pageSize

	// Fetch paginated results with sorting
	fetchQuery := transaction.Model(&models.DBMessage{}).
		Where("tenant_id = ?", tenantId)

	// Apply same status filter to fetch query
	if status != "" && (status == "pending" || status == "approved" || status == "rejected") {
		fetchQuery = fetchQuery.Where("status = ?", status)
	}

	// Execute query with pagination and ordering
	// Order by: pending first, then by most recent request
	if err := fetchQuery.
		Order("CASE WHEN status = 'pending' THEN 0 ELSE 1 END"). // Pending items first
		Order("request_at DESC").                                // Most recent first
		Limit(pageSize).
		Offset(offset).
		Find(&messages).Error; err != nil {
		log.Printf("Error fetching messages: %v", err)
		return nil, 0, fmt.Errorf("error fetching messages: %v", err)
	}

	// Commit transaction
	if err := transaction.Commit().Error; err != nil {
		return nil, 0, err
	}

	return messages, totalCount, nil
}

func (tm *TenantMessageRepository) ApproveMessage(tenantId uuid.UUID, messageId uuid.UUID) error {
	transaction := tm.DB.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	var message *models.DBMessage
	update := tm.DB.Model(&models.DBMessage{}).Where("tenant_id = ? ", tenantId).Where("id = ? ", messageId).First(&message)
	if update.Error != nil {
		return update.Error
	}
	if message.Action {
		return fmt.Errorf("error while updating the status as the this query is already been updated and the status is %s", message.Status)
	}
	updateQuery := tm.DB.Model(&models.DBMessage{}).Where("tenant_id = ? ", tenantId).Where("id = ? ", messageId).Updates(map[string]interface{}{
		"action": true,
		"status": "approved"})
	if updateQuery.Error != nil {
		return updateQuery.Error
	}
	transaction.Commit()
	return nil
}

func (tm *TenantMessageRepository) RejectMessage(tenantId uuid.UUID, messageId uuid.UUID) error {
	transaction := tm.DB.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	var message *models.DBMessage
	update := tm.DB.Model(&models.DBMessage{}).Where("tenant_id = ? ", tenantId).Where("id = ? ", messageId).First(&message)
	if update.Error != nil {
		return update.Error
	}
	if message.Action {
		return fmt.Errorf("error while updating the status as the this query is already been updated and the status is %s", message.Status)
	}
	updateQuery := tm.DB.Model(&models.DBMessage{}).Where("tenant_id = ? ", tenantId).Where("id = ? ", messageId).Updates(map[string]interface{}{
		"action": true,
		"status": "rejected"})
	if updateQuery.Error != nil {
		return updateQuery.Error
	}
	transaction.Commit()
	return nil
}
