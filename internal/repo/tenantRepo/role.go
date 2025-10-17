package tenantrepo

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

type TenantRoleRepositoryInterface interface {
	ListRoles(tenantId uuid.UUID, page int, pageSize int, status string, roleType string) ([]*models.DBRoles, int64, error)
}

type TenantRoleRepository struct {
	DB *gorm.DB
}

func NewTenantRoleRepository(db *gorm.DB) (TenantRoleRepositoryInterface, error) {
	return &TenantRoleRepository{DB: db}, nil
}

func (tr *TenantRoleRepository) ListRoles(tenantId uuid.UUID, page int, pageSize int, status string, roleType string) ([]*models.DBRoles, int64, error) {
	var totalCount int64

	var roles []*models.DBRoles
	transaction := tr.DB.Begin()
	if transaction.Error != nil {
		return nil, 0, transaction.Error
	}
	var is_active bool
	if status == "enabled" {
		is_active = true
	} else {
		is_active = false
	}
	defer transaction.Rollback()
	if err := tr.DB.Model(&models.DBRoles{}).Where("tenant_id = ?", tenantId).Where("status = ?", is_active).Where("role_type = ? ", roleType).Count(&totalCount).Error; err != nil {
		log.Printf("Error counting the tenant-roles: %v", err)
		return nil, 0, fmt.Errorf("error counting the roles present for this tenant: %v", err)
	}
	log.Printf("Total tokens found: %d", totalCount)


	offset := (page - 1) * pageSize

	if err := tr.DB.Model(&models.DBRoles{}).
		Where("tenant_id = ?", tenantId).Where("status = ?", is_active).Where("role_type = ? ", roleType).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&roles).Error; err != nil {
		log.Printf("Error fetching paginated roles for the tenant: %v", err)
		return nil, 0, fmt.Errorf("error fetching roles for this tenant: %w", err)
	}
log.Printf("Successfully fetched %d roles for this tenant (page %d, pageSize %d, total %d)",
		len(roles), page, pageSize, totalCount)

	return roles, totalCount, nil
}
