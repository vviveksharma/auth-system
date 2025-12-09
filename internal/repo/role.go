package repo

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

type RoleRepositoryInterface interface {
	CreateRole(req *models.DBRoles) error
	GetAllRoles(roleTypeFlag string, tenantId uuid.UUID, page, pageSize int) ([]*models.DBRoles, int64, error)
	FindRoleId(roleName string) (roleId uuid.UUID, err error)
	DeleteRole(roleId uuid.UUID) error
	GetRolesDetails(conditions *models.DBRoles) (resp *models.DBRoles, err error)
	ChangeStatus(flag bool, roleId uuid.UUID) error
	GetRolesByTenant(tenantId uuid.UUID, roleType string) ([]*models.DBRoles, error)
	UpdateRoleDetails(req models.DBRoles, tenantId uuid.UUID) error
}

type RoleRepository struct {
	DB *gorm.DB
}

func NewRoleRepository(db *gorm.DB) (RoleRepositoryInterface, error) {
	return &RoleRepository{DB: db}, nil
}

func (r *RoleRepository) GetAllRoles(roleTypeFlag string, tenantId uuid.UUID, page, pageSize int) ([]*models.DBRoles, int64, error) {
	fmt.Println("Starting GetAllRoles transaction with pagination")

	var totalCount int64
	var roleDetails []*models.DBRoles

	if roleTypeFlag == "default" {
		tenant := models.GetSystemTenantId()
		tenantId = uuid.MustParse(tenant)
	}

	baseQuery := r.DB.Model(&models.DBRoles{}).Where("tenant_id = ?", tenantId)
	if err := baseQuery.Count(&totalCount).Error; err != nil {
		fmt.Printf("Error counting roles in GetAllRoles: %v\n", err)
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	query := r.DB.Where("tenant_id = ?", tenantId)
	if roleTypeFlag != "" && roleTypeFlag != "all" {
		query = query.Where("role_type = ?", roleTypeFlag)
	}

	if err := query.Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&roleDetails).Error; err != nil {
		fmt.Printf("Error fetching paginated roles in GetAllRoles: %v\n", err)
		return nil, 0, err
	}

	fmt.Printf("Successfully fetched %d roles (page %d, pageSize %d, total %d)\n",
		len(roleDetails), page, pageSize, totalCount)

	return roleDetails, totalCount, nil
}

func (r *RoleRepository) FindRoleId(roleName string) (roleId uuid.UUID, err error) {
	fmt.Printf("Starting FindRoleId transaction for roleName: %s\n", roleName)
	transaction := r.DB.Begin()
	if transaction.Error != nil {
		fmt.Printf("Failed to begin transaction in FindRoleId: %v\n", transaction.Error)
		return uuid.Nil, transaction.Error
	}
	defer func() {
		fmt.Println("Rolling back FindRoleId transaction")
		transaction.Rollback()
	}()
	var roles models.DBRoles
	result := transaction.Where("role = ?", roleName).First(&roles)

	if result.Error != nil {
		return uuid.Nil, result.Error
	}
	fmt.Printf("Successfully found roleId %s for roleName '%s'\n", roles.RoleId, roleName)
	return roles.RoleId, nil
}

func (r *RoleRepository) GetRolesDetails(conditions *models.DBRoles) (resp *models.DBRoles, err error) {
	transaction := r.DB.Begin()
	if transaction.Error != nil {
		return nil, transaction.Error
	}
	defer transaction.Rollback()
	rerr := transaction.First(&resp, &conditions)
	if rerr.Error != nil {
		return nil, rerr.Error
	}
	return resp, nil
}

func (r *RoleRepository) CreateRole(req *models.DBRoles) error {
	transaction := r.DB.Begin()
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

func (r *RoleRepository) DeleteRole(roleId uuid.UUID) error {
	transaction := r.DB.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	delete := transaction.Model(models.DBRoles{}).Where("role_id = ?", roleId).Delete(models.DBRoles{
		RoleId: roleId,
	})
	if delete.Error != nil {
		return delete.Error
	}
	return nil
}

func (r *RoleRepository) ChangeStatus(flag bool, roleId uuid.UUID) error {
	transaction := r.DB.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	if !flag { // disable the role
		update := transaction.Model(&models.DBRoles{}).Where("role_id = ?", roleId).Updates(map[string]interface{}{
			"status": false,
		})
		if update.Error != nil {
			return update.Error
		}
	} else { // enable the role
		update := transaction.Model(&models.DBRoles{}).Where("role_id = ?", roleId).Updates(map[string]interface{}{
			"status": true,
		})
		if update.Error != nil {
			return update.Error
		}
	}
	transaction.Commit()
	return nil
}

func (r *RoleRepository) GetRolesByTenant(tenantId uuid.UUID, roleType string) ([]*models.DBRoles, error) {
	var roles []*models.DBRoles
	query := r.DB.Where("tenant_id = ?", tenantId)

	if roleType != "" && roleType != "all" {
		query = query.Where("role_type = ?", roleType)
	}

	err := query.Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *RoleRepository) UpdateRoleDetails(req models.DBRoles, tenantId uuid.UUID) error {
	fmt.Printf("UpdateRoleDetails called for roleId: %s, tenantId: %s\n", req.RoleId, tenantId)

	transaction := r.DB.Begin()
	if transaction.Error != nil {
		fmt.Printf("Error starting transaction: %v\n", transaction.Error)
		return transaction.Error
	}
	defer transaction.Rollback()

	var existingRole models.DBRoles
	if err := transaction.Where("role_id = ? AND tenant_id = ?", req.RoleId, tenantId).First(&existingRole).Error; err != nil {
		fmt.Printf("Role not found or doesn't belong to tenant: %v\n", err)
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("role not found with id %s for tenant %s", req.RoleId, tenantId)
		}
		return fmt.Errorf("error fetching role: %w", err)
	}

	updates := make(map[string]interface{})
	if req.Role != "" {
		updates["role"] = req.Role
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	updates["status"] = req.Status

	result := transaction.Model(&models.DBRoles{}).
		Where("role_id = ? AND tenant_id = ?", req.RoleId, tenantId).
		Updates(updates)

	if result.Error != nil {
		fmt.Printf("Error updating role: %v\n", result.Error)
		return fmt.Errorf("error updating role: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		fmt.Printf("No rows updated for roleId: %s\n", req.RoleId)
		return fmt.Errorf("no role found to update with id %s", req.RoleId)
	}

	if err := transaction.Commit().Error; err != nil {
		fmt.Printf("Error committing transaction: %v\n", err)
		return fmt.Errorf("error committing transaction: %w", err)
	}

	fmt.Printf("Successfully updated role %s (rows affected: %d)\n", req.RoleId, result.RowsAffected)
	return nil
}
