package repo

import (
	"log"
	"slices"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

type RouteRoleRepositoryInterface interface {
	Create(req *models.DBRouteRole) error
	FindByRoleId(roleId uuid.UUID) (bool, error)
	UpdateRouteRole(roleId string, route string) error
	DeleteAndUpdateRole(roleId string, addroutes []string, removeroutes []string) error
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

func (rr *RouteRoleRepository) DeleteAndUpdateRole(roleId string, addroutes []string, removeroutes []string) error {
	transaction := rr.DB.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	//find the roleDetails as per the roleId
	var RoleRouteDetails models.DBRouteRole
	roleRoute := transaction.Where("role_id = ?", uuid.MustParse(roleId)).First(&RoleRouteDetails)
	if roleRoute.Error != nil {
		return roleRoute.Error
	}
	filtered := []string{}
	if len(removeroutes) > 0 {
		removeSet := map[string]struct{}{}
		for _, r := range removeroutes {
			removeSet[r] = struct{}{}
		}

		var filteredRoutes []string
		for _, r := range RoleRouteDetails.Route {
			if _, found := removeSet[r]; !found {
				filteredRoutes = append(filteredRoutes, r)
			}
		}
		log.Println("the routes after remove: ", filteredRoutes)
		filtered = filteredRoutes
	}

	// Add new routes
	filtered = append(filtered, addroutes...)
	log.Println("the routes after remove: ", filtered)
	

	// Remove duplicates if any
	filtered = uniqueStrings(filtered)
	RoleRouteDetails.Route = RoleRouteDetails.Route[:0]
	update := transaction.Model(&models.DBRouteRole{}).Where("role_id = ?", uuid.MustParse(roleId)).Updates(map[string]interface{}{
		"route": pq.StringArray(filtered),
	})
	if update.Error != nil {
		return update.Error
	}
	transaction.Commit()
	return nil
}

func uniqueStrings(input []string) []string {
	seen := make(map[string]struct{})
	result := []string{}
	for _, item := range input {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
