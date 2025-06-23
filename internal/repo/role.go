package repo

import (
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

type RoleRepositoryInterface interface {
	GetAllRoles() ([]*models.DBRoles, error)
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
