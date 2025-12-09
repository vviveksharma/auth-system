package tenantrepo

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

type TenantUserRepositoryInterface interface {
	ListUsers(page int, pageSize int, tenantId uuid.UUID, status string) ([]*models.DBUser, int64, error)
}

type TenantUserRepository struct {
	DB *gorm.DB
}

func NewTenantUserRepository(db *gorm.DB) (TenantUserRepositoryInterface, error) {
	return &TenantUserRepository{DB: db}, nil
}

func (tu *TenantUserRepository) ListUsers(page int, pageSize int, tenantId uuid.UUID, status string) ([]*models.DBUser, int64, error) {
	log.Printf("🔍 [TENANT REPO] ListUsers called with: page=%d, pageSize=%d, tenantId=%s, status=%s",
		page, pageSize, tenantId.String(), status)

	var totalCount int64
	var users []*models.DBUser

	// Build count query
	countQuery := tu.DB.Model(&models.DBUser{}).Where("tenant_id = ?", tenantId)

	// Apply status filter only if not "all"
	if status != "all" {
		var is_active bool
		switch status {
		case "active":
			is_active = true
		case "inactive":
			is_active = false
		}
		log.Printf("🔍 [TENANT REPO] Converted status '%s' to is_active=%v", status, is_active)
		countQuery = countQuery.Where("status = ?", is_active)
	} else {
		log.Printf("🔍 [TENANT REPO] Status is 'all', fetching all users regardless of status")
	}

	if err := countQuery.Count(&totalCount).Error; err != nil {
		log.Printf("❌ [TENANT REPO] Error counting users: %v", err)
		return nil, 0, fmt.Errorf("error counting the users present for this tenant: %v", err)
	}
	log.Printf("✅ [TENANT REPO] Total users found: %d", totalCount)

	offset := (page - 1) * pageSize
	log.Printf("🔍 [TENANT REPO] Calculated offset: %d (page=%d, pageSize=%d)", offset, page, pageSize)

	// Build fetch query
	fetchQuery := tu.DB.Model(&models.DBUser{}).
		Where("tenant_id = ?", tenantId)

	// Apply status filter only if not "all"
	if status != "all" {
		var is_active bool
		if status == "enabled" {
			is_active = true
		} else {
			is_active = false
		}
		fetchQuery = fetchQuery.Where("status = ?", is_active)
	}

	fetchQuery = fetchQuery.Order("created_at DESC").
		Limit(pageSize).
		Offset(offset)

	if err := fetchQuery.Find(&users).Error; err != nil {
		log.Printf("❌ [TENANT REPO] Error fetching paginated users: %v", err)
		return nil, 0, fmt.Errorf("error fetching users for this tenant: %w", err)
	}

	log.Printf("✅ [TENANT REPO] Successfully fetched %d users (page %d, pageSize %d, total %d)",
		len(users), page, pageSize, totalCount)

	for i, user := range users {
		log.Printf("   [TENANT REPO] User %d: ID=%s, Email=%s, Name=%s, Status=%v",
			i+1, user.Id.String(), user.Email, user.Name, user.Status)
	}

	return users, totalCount, nil
}
