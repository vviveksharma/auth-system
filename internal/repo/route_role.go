package repo

import (
	"slices"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

type RouteRoleRepositoryInterface interface {
	Create(req *models.DBRouteRole) error
	FindByRoleId(roleId uuid.UUID) (bool, error)
	UpdateRouteRole(roleId string, route string) error
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

func (rr *RouteRoleRepository) FindByRoleId(roleId uuid.UUID) (bool, error) {
	transaction := rr.DB.Begin()
	if transaction.Error != nil {
		return false, transaction.Error
	}
	defer transaction.Rollback()
	var RoleRouteDetails models.DBRouteRole
	err := rr.DB.Where("role_id = ?", roleId).First(&RoleRouteDetails).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (rr *RouteRoleRepository) UpdateRouteRole(roleId string, route string) error {
	transaction := rr.DB.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	var RoleRouteDetails models.DBRouteRole
	roleRouteDetails := transaction.Model(&models.DBRouteRole{}).Where("role_id = ?", roleId).Find(&RoleRouteDetails)
	if roleRouteDetails.Error != nil {
		return roleRouteDetails.Error
	}
	if slices.Contains(RoleRouteDetails.Route, route) {
		return nil
	}

	// Append the new route and update
	RoleRouteDetails.Route = append(RoleRouteDetails.Route, route)
	update := rr.DB.Model(&models.DBRouteRole{}).Where("role_id = ?", roleId).Update("route", RoleRouteDetails.Route)
	if update.Error != nil {
		return update.Error
	}

	return nil

}
