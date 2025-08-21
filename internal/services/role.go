package services

import (
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/vviveksharma/auth/db"
	"github.com/vviveksharma/auth/internal/models"
	"github.com/vviveksharma/auth/internal/repo"
	dbmodels "github.com/vviveksharma/auth/models"
)

type RoleService interface {
	CreateCustomRole(req *models.CreateCustomRole) (resp *models.CreateCustomRoleResponse, err error)
	ListAllRoles(typeFlag string) (response []*models.ListAllRolesResponse, err error)
	VerifyRole(req *models.VerifyRoleRequest) (response *models.VerifyRoleResponse, err error)
	UpdateRolePermission(req *models.UpdateRolePermissions) (resp *models.CreateCustomRoleResponse, err error)
}

type Role struct {
	RoleRepo      repo.RoleRepositoryInterface
	RoleRouteRepo repo.RouteRoleRepositoryInterface
}

func NewRoleService() (RoleService, error) {
	ser := &Role{}
	err := ser.SetupRepo()
	if err != nil {
		return nil, err
	}
	return ser, nil
}

func (r *Role) SetupRepo() error {
	var err error
	role, err := repo.NewRoleRepository(db.DB)
	if err != nil {
		return err
	}
	r.RoleRepo = role
	routeRole, err := repo.NewRouteRoleRepository(db.DB)
	if err != nil {
		return err
	}
	r.RoleRouteRepo = routeRole
	return nil
}

func (r *Role) ListAllRoles(typeFlag string) (response []*models.ListAllRolesResponse, err error) {
	roles, err := r.RoleRepo.GetAllRoles()
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while listing all the roles: " + err.Error(),
		}
	}
	for _, i := range roles {
		if i.RoleType != typeFlag {
			continue
		}
		var res models.ListAllRolesResponse
		res.Name = i.Role
		response = append(response, &res)
	}
	return response, nil
}

func (r *Role) VerifyRole(req *models.VerifyRoleRequest) (response *models.VerifyRoleResponse, err error) {
	roleDetails, err := r.RoleRepo.FindRoleId(req.RoleName)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "No role with this name exists",
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "error while finding the role:" + err.Error(),
			}
		}
	}
	log.Info("the roleId: ", roleDetails)
	if roleDetails != uuid.MustParse(req.RoleId) {
		return nil, &dbmodels.ServiceResponse{
			Code:    404,
			Message: "Role name and ID do not match. Please provide the correct role ID and role name.",
		}
	}
	return &models.VerifyRoleResponse{
		Message: true,
	}, nil
}

func (r *Role) CreateCustomRole(req *models.CreateCustomRole) (resp *models.CreateCustomRoleResponse, err error) {
	// First fetching the role id
	roleId, err := r.RoleRepo.FindRoleId(req.RoleName)
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "An unexpected error occurred while searching for the role '" + req.RoleName + "': " + err.Error(),
			}
		} else {
			newRoleId := uuid.New()
			err := r.RoleRepo.CreateRole(&dbmodels.DBRoles{
				Role:     req.RoleName,
				RoleId:   newRoleId,
				RoleType: "custom",
			})
			if err != nil {
				return nil, &dbmodels.ServiceResponse{
					Code:    500,
					Message: "An unexpected error occurred while creating the role: " + err.Error(),
				}
			}
			roleId = newRoleId
		}

	}

	rresp, err := r.RoleRouteRepo.FindByRoleId(roleId)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "An unexpected error occurred while searching for the role-route association: " + err.Error(),
		}
	}
	if rresp {
		// RoleId entry present, updating the routes (adding them)
		for _, route := range req.Routes {
			updateRouteErr := r.RoleRouteRepo.UpdateRouteRole(roleId.String(), route)
			if updateRouteErr != nil {
				return nil, &dbmodels.ServiceResponse{
					Code:    500,
					Message: "An unexpected error occurred while updating the route for the existing role ID: " + updateRouteErr.Error(),
				}
			}
		}
	} else {
		err := r.RoleRouteRepo.Create(&dbmodels.DBRouteRole{
			TenantId: uuid.MustParse("dae760ab-0a7f-4cbd-8603-def85ad8e430"),
			RoleId:   roleId,
			Routes:    req.Routes,
		})
		if err != nil {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "An unexpected error occurred while creating a new role-route entry: " + err.Error(),
			}
		}
	}
	return &models.CreateCustomRoleResponse{
		Message: "role with " + req.RoleName + " created successfully.",
	}, nil
}

func (r *Role) UpdateRolePermission(req *models.UpdateRolePermissions) (resp *models.CreateCustomRoleResponse, err error) {
	roleDetails, err := r.RoleRepo.FindByName(req.RoleName)
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "An unexpected error occurred while searching for the role '" + req.RoleName + "': " + err.Error(),
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "The requested role name '" + req.RoleName + "' does not exist. Please verify the role name and try again.",
			}
		}
	}
	if roleDetails.RoleType == "default" {
		return nil, &dbmodels.ServiceResponse{
			Code:    409,
			Message: "The default role cannot be modified. Please create a custom role if you need to make changes.",
		}
	}
	err = r.RoleRouteRepo.DeleteAndUpdateRole(roleDetails.RoleId.String(), req.AddPermisions, req.RemovePermissions)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while updating the routes in the: " + err.Error(),
		}
	}
	return &models.CreateCustomRoleResponse{
		Message: "The role permissions for '" + req.RoleName + "' have been updated successfully.",
	}, nil
}
