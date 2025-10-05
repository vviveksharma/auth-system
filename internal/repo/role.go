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
	FindByName(roleName string) (*models.DBRoles, error)
	DeleteRole(roleId uuid.UUID) error
	GetRolesDetails(conditions *models.DBRoles) (resp *models.DBRoles, err error)
	ChangeStatus(flag bool, roleId uuid.UUID) error
	GetRoleByName(roleName string, tenantId uuid.UUID) (*models.DBRoles, error)
	GetRolesByTenant(tenantId uuid.UUID, roleType string) ([]*models.DBRoles, error)
	GetRoleUsageCount(roleId uuid.UUID, tenantId string) (int64, error)
	IsSystemRole(roleId uuid.UUID) (bool, error)
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

func (r *RoleRepository) FindByName(roleName string) (resp *models.DBRoles, err error) {
	transaction := r.DB.Begin()
	if transaction.Error != nil {
		return nil, transaction.Error
	}
	defer transaction.Rollback()
	rr := transaction.Where("role = ?", roleName).First(&resp)
	if rr.Error != nil {
		return nil, rr.Error
	}
	return resp, nil
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

// **NEW: Additional methods for seeding**
func (r *RoleRepository) GetRoleByName(roleName string, tenantId uuid.UUID) (*models.DBRoles, error) {
	var role models.DBRoles
	err := r.DB.Where("role = ? AND tenant_id = ?", roleName, tenantId).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
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

func (r *RoleRepository) GetRoleUsageCount(roleId uuid.UUID, tenantId string) (int64, error) {
	var count int64
	// This assumes users have roles in a roles array - adjust based on your user model
	err := r.DB.Model(&models.DBUser{}).
		Where("tenant_id = ? AND ? = ANY(roles)", tenantId, roleId.String()).
		Count(&count).Error
	return count, err
}

func (r *RoleRepository) IsSystemRole(roleId uuid.UUID) (bool, error) {
	transaction := r.DB.Begin()
	if transaction.Error != nil {
		return false, transaction.Error
	}
	defer transaction.Rollback()
	var roleDetails models.DBRoles
	err := transaction.Where("role_id = ? ", roleId).Find(&roleDetails)
	if err.Error != nil {
		return false, err.Error
	}
	if roleDetails.RoleType != "default" {
		return false, nil
	}
	return true, nil
}
