package repo

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

type RoleRepositoryInterface interface {
	GetAllRoles() ([]*models.DBRoles, error)
	FindRoleId(roleName string) (roleId uuid.UUID, err error)
}

type RoleRepository struct {
	DB *gorm.DB
}

func NewRoleRepository(db *gorm.DB) (RoleRepositoryInterface, error) {
	return &RoleRepository{DB: db}, nil
}

func (r *RoleRepository) GetAllRoles() ([]*models.DBRoles, error) {
	transaction := r.DB.Begin()
	if transaction.Error != nil {
		return nil, transaction.Error
	}
	defer transaction.Rollback()
	roleDetails := []*models.DBRoles{}
	roles := transaction.Find(&roleDetails)
	if roles.Error != nil {
		return nil, roles.Error
	}
	return roleDetails, nil
}

func (r *RoleRepository) FindRoleId(roleName string) (roleId uuid.UUID, err error) {
	transaction := r.DB.Begin()
	if transaction.Error != nil {
		return uuid.Nil, transaction.Error
	}
	defer transaction.Rollback()
	var roles models.DBRoles
	result := transaction.Where("role = ?", roleName).First(&roles)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return uuid.Nil, fmt.Errorf("role '%s' not found", roleName)
	}
	if result.Error != nil {
		return uuid.Nil, result.Error
	}
	return roles.Id, nil
}
