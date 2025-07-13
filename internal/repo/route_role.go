package repo

import (
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

type RouteRoleRepositoryInterface interface {
	Create(req models.DBRouteRole) error
}

type RouteRoleRepository struct {
	DB *gorm.DB
}

func NewRouteRoleRepository(db *gorm.DB) (RouteRoleRepositoryInterface, error) {
	return &RouteRoleRepository{DB: db}, nil
}

func (rr *RouteRoleRepository) Create(req models.DBRouteRole) error {
	transaction := rr.DB.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	err := rr.DB.Create(&req)
	if err.Error != nil {
		return err.Error
	}
	return nil
}
