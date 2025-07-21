package repo

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

type RoleRepositoryInterface interface {
	CreateRole(req *models.DBRoles) error
	GetAllRoles() ([]*models.DBRoles, error)
	FindRoleId(roleName string) (roleId uuid.UUID, err error)
	FindByName(roleName string) (*models.DBRoles, error)
}

type RoleRepository struct {
	DB *gorm.DB
}

func NewRoleRepository(db *gorm.DB) (RoleRepositoryInterface, error) {
	return &RoleRepository{DB: db}, nil
}

func (r *RoleRepository) GetAllRoles() ([]*models.DBRoles, error) {
	fmt.Println("Starting GetAllRoles transaction")
	transaction := r.DB.Begin()
	if transaction.Error != nil {
		fmt.Printf("Failed to begin transaction in GetAllRoles: %v\n", transaction.Error)
		return nil, transaction.Error
	}
	defer func() {
		fmt.Println("Rolling back GetAllRoles transaction")
		transaction.Rollback()
	}()
	roleDetails := []*models.DBRoles{}
	roles := transaction.Find(&roleDetails)
	if roles.Error != nil {
		fmt.Printf("Error fetching roles in GetAllRoles: %v\n", roles.Error)
		return nil, roles.Error
	}
	fmt.Printf("Successfully fetched %d roles\n", len(roleDetails))
	return roleDetails, nil
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
