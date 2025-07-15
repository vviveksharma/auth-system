package repo

import (
	"errors"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

type RouteRoleRepositoryInterface interface {
	Create(req *models.DBRouteRole) error
	FindByRoute(route string, roleId uuid.UUID) (bool, error)
}

type RouteRoleRepository struct {
	DB *gorm.DB
}

func NewRouteRoleRepository(db *gorm.DB) (RouteRoleRepositoryInterface, error) {
	return &RouteRoleRepository{DB: db}, nil
}

func (rr *RouteRoleRepository) Create(req *models.DBRouteRole) error {
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

func (rr *RouteRoleRepository) FindByRoute(route string, roleId uuid.UUID) (bool, error) {
	transaction := rr.DB.Begin()
	if transaction.Error != nil {
		return false, transaction.Error
	}
	defer transaction.Rollback()
	var RoleRouteDetails models.DBRouteRole
	rrErr := transaction.Model(&models.DBRouteRole{}).Where("route = ?", route).Find(&RoleRouteDetails)
	if rrErr.Error != nil {
		return false, rrErr.Error
	}
	for _, role := range RoleRouteDetails.RoleId {
		if uuid.MustParse(role) == roleId {
			return true, nil
		}
	}
	return false, errors.New("the route associated with role not found")
}
