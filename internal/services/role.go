package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/vviveksharma/auth/db"
	"github.com/vviveksharma/auth/internal/models"
	"github.com/vviveksharma/auth/internal/pagination"
	"github.com/vviveksharma/auth/internal/repo"
	dbmodels "github.com/vviveksharma/auth/models"
)

type RoleService interface {
	CreateCustomRole(req *models.CreateCustomRole, ctx context.Context) (resp *models.CreateCustomRoleResponse, err error)
	VerifyRole(req *models.VerifyRoleRequest) (response *models.VerifyRoleResponse, err error)
	ListRoles(roleTypeFlag string, page int, pageSize int, ctx context.Context) (*pagination.PaginatedResponse[*models.ListAllRolesResponse], error)
	UpdateRolePermission(req *models.UpdateRolePermissions, roleId uuid.UUID, ctx context.Context) (resp *models.CreateCustomRoleResponse, err error)
	DeleteRole(roleId uuid.UUID, ctx context.Context) (*models.DeleteRoleResponse, error)
	EnableRole(roleId uuid.UUID, ctx context.Context) (*models.EnableRoleResponse, error)
	DisableRole(roleId uuid.UUID, ctx context.Context) (*models.DisableRoleResponse, error)
	GetRouteDetails(roleId uuid.UUID, ctx context.Context) (*models.GetRouteDetailsResponse, error)
}

type Role struct {
	RoleRepo      repo.RoleRepositoryInterface
	RoleRouteRepo repo.RouteRoleRepositoryInterface
	SharedRepo    repo.SharedRepoInterface
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
	sharedRepo, err := repo.NewSharedRepository(db.DB)
	if err != nil {
		return err
	}
	r.SharedRepo = sharedRepo
	return nil
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

func (r *Role) CreateCustomRole(req *models.CreateCustomRole, ctx context.Context) (resp *models.CreateCustomRoleResponse, err error) {
	// First fetching the role id
	tenantId := ctx.Value("tenant_id").(string)
	roleDetails, err := r.RoleRepo.GetRoleByName(req.Name, uuid.MustParse(tenantId))
	if err != nil && err.Error() != "record not found" {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while fetching the roleDetails from the db: " + err.Error(),
		}
	}
	if roleDetails != nil && roleDetails.Role == req.Name {
		return nil, &dbmodels.ServiceResponse{
			Code:    409,
			Message: "role with this name alredy exist please go ahead with the updating the role",
		}
	}
	// Create the role and role Route mapping
	err = r.SharedRepo.CreateCustomRole(req, uuid.MustParse(tenantId))
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while creating the custom role: " + err.Error(),
		}
	}
	return &models.CreateCustomRoleResponse{
		Message: fmt.Sprintf("Role with %s created successfully: ", req.Name),
	}, nil
}

func (r *Role) UpdateRolePermission(req *models.UpdateRolePermissions, roleId uuid.UUID, ctx context.Context) (resp *models.CreateCustomRoleResponse, err error) {
	tenantId := ctx.Value("tenant_id").(string)
	roleDetails, err := r.RoleRepo.GetRolesDetails(&dbmodels.DBRoles{
		RoleId:   roleId,
		TenantId: uuid.MustParse(tenantId),
	})
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
	err = r.SharedRepo.UpdateCustomRole(roleDetails.RoleId, roleDetails.TenantId, req.AddPermisions, req.RemovePermissions)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while updating the role-route mapping: " + err.Error(),
		}
	}

	return &models.CreateCustomRoleResponse{
		Message: "The role permissions for '" + req.RoleName + "' have been updated successfully.",
	}, nil
}

func (r *Role) DeleteRole(roleId uuid.UUID, ctx context.Context) (*models.DeleteRoleResponse, error) {
	tenantId := ctx.Value("tenant_id").(string)
	roleDetails, err := r.RoleRepo.GetRolesDetails(&dbmodels.DBRoles{
		RoleId:   roleId,
		TenantId: uuid.MustParse(tenantId),
	})
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "No role found with the provided role ID. Please verify the role ID and try again.",
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "An unexpected error occurred while retrieving role details from the database: " + err.Error(),
			}
		}
	}
	if roleDetails.RoleType == "default" {
		return nil, &dbmodels.ServiceResponse{
			Code:    404,
			Message: "System-level roles cannot be deleted. Only custom roles created by you can be removed.",
		}
	}
	err = r.SharedRepo.DeleteCustomRole(roleId, uuid.MustParse(tenantId))
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "An error occurred while attempting to delete the role: " + err.Error() + ". Please contact support if the issue persists.",
		}
	}
	return &models.DeleteRoleResponse{
		Message: fmt.Sprintf("The role with ID '%s' has been deleted successfully.", roleDetails.RoleId.String()),
	}, nil
}

func (r *Role) EnableRole(roleId uuid.UUID, ctx context.Context) (*models.EnableRoleResponse, error) {
	tenantId := ctx.Value("tenant_id").(string)
	// fetch role details
	roleDetails, err := r.RoleRepo.GetRolesDetails(&dbmodels.DBRoles{
		TenantId: uuid.MustParse(tenantId),
		RoleId:   roleId,
	})
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "No role found with the provided role ID. Please verify the role ID and try again.",
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "An unexpected error occurred while retrieving role details from the database: " + err.Error(),
			}
		}
	}
	if r.isSystem(roleDetails) {
		return nil, &dbmodels.ServiceResponse{
			Code:    403,
			Message: "System-level roles cannot be modified. This operation is restricted to custom roles only.",
		}
	}
	if roleDetails.Status {
		return nil, &dbmodels.ServiceResponse{
			Code:    409,
			Message: "The role is already enabled and active. No action is required.",
		}
	}
	err = r.RoleRepo.ChangeStatus(true, roleId)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "Failed to enable the role due to a database error: " + err.Error() + ". Please try again or contact support if the issue persists.",
		}
	}
	return &models.EnableRoleResponse{
		Message: fmt.Sprintf("Role with ID '%s' has been enabled successfully.", roleId.String()),
	}, nil
}

func (r *Role) DisableRole(roleId uuid.UUID, ctx context.Context) (*models.DisableRoleResponse, error) {
	tenantId := ctx.Value("tenant_id").(string)
	// fetch role details
	roleDetails, err := r.RoleRepo.GetRolesDetails(&dbmodels.DBRoles{
		TenantId: uuid.MustParse(tenantId),
		RoleId:   roleId,
	})
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "No role found with the provided role ID. Please verify the role ID and try again.",
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "An unexpected error occurred while retrieving role details from the database: " + err.Error(),
			}
		}
	}
	if r.isSystem(roleDetails) {
		return nil, &dbmodels.ServiceResponse{
			Code:    403,
			Message: "System-level roles cannot be modified. This operation is restricted to custom roles only.",
		}
	}
	if roleDetails.Status {
		return nil, &dbmodels.ServiceResponse{
			Code:    409,
			Message: "The role is already disabled and active. No action is required.",
		}
	}
	err = r.RoleRepo.ChangeStatus(true, roleId)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "Failed to disable the role due to a database error: " + err.Error() + ". Please try again or contact support if the issue persists.",
		}
	}
	return &models.DisableRoleResponse{
		Message: fmt.Sprintf("Role with ID '%s' has been disbaled successfully.", roleId.String()),
	}, nil
}

func (r *Role) GetRouteDetails(roleId uuid.UUID, ctx context.Context) (*models.GetRouteDetailsResponse, error) {
	tenantId := ctx.Value("tenant_id").(string)
	// check the roleDetails value and then select the tenantId
	isSystemRole, err := r.RoleRepo.IsSystemRole(roleId)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "No role found with the provided role ID. Please verify the role ID and try again.",
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "An unexpected error occurred while retrieving role details from the database while checking the isSystemRole: " + err.Error(),
			}
		}
	}
	if isSystemRole {
		tenantId = dbmodels.GetSystemTenantId()
	}
	// fetch role details
	roleDetails, err := r.RoleRepo.GetRolesDetails(&dbmodels.DBRoles{
		TenantId: uuid.MustParse(tenantId),
		RoleId:   roleId,
	})
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "No role found with the provided role ID. Please verify the role ID and try again.",
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "An unexpected error occurred while retrieving role details from the database: " + err.Error(),
			}
		}
	}
	roleRoute, err := r.RoleRouteRepo.GetRoleRouteMapping(roleDetails.RoleId.String())
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "An unexpected error occurred while retrieving role details from the database: " + err.Error(),
		}
	}

	fmt.Println("the roleRoute response: ", roleRoute.Permissions)

	permissions, err := models.ConvertDBData(roleRoute.Permissions)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while converting the permissions to the JSON: " + err.Error(),
		}
	}

	fmt.Println("the permissions: ", permissions.Permissions)

	classifiedPermissions := models.ClassifyPermissionsByMethod(permissions.Permissions)

	resp, err := json.MarshalIndent(classifiedPermissions, "", "  ")
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "erroir while converting to the JSON: " + err.Error(),
		}
	}

	return &models.GetRouteDetailsResponse{
		Routes:      classifiedPermissions,
		RoutesJSON:  string(resp),
		RoleInfo:    permissions.RoleInfo,
		ProcessedAt: time.Now(),
	}, nil
}

func (r *Role) isSystem(roleDetails *dbmodels.DBRoles) bool {
	return roleDetails.RoleType == "default"
}

func (r *Role) ListRoles(roleTypeFlag string, page int, pageSize int, ctx context.Context) (*pagination.PaginatedResponse[*models.ListAllRolesResponse], error) {
	tenantId := ctx.Value("tenant_id").(string)
	fmt.Println("the roletypeflag: ", roleTypeFlag)
	roleDetails, totalCount, err := r.RoleRepo.GetAllRoles(roleTypeFlag, uuid.MustParse(tenantId), page, pageSize)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while listing roles: " + err.Error(),
		}
	}
	var responseRoles []*models.ListAllRolesResponse
	for _, role := range roleDetails {
		roleResponse := &models.ListAllRolesResponse{
			Name:     role.Role,
			RoleId:   role.RoleId,
			RoleType: role.RoleType,
			Status:   role.Status,
			TenantId: role.TenantId,
		}

		routes, err := r.RoleRouteRepo.GetRoleRouteMapping(role.RoleId.String())
		if err != nil {
			fmt.Printf("Warning: Could not fetch routes for role %s: %v\n", role.RoleId.String(), err)
			roleResponse.Routes = []string{}
		} else {
			roleResponse.Routes = routes.Routes
		}

		responseRoles = append(responseRoles, roleResponse)
	}
	totalPages := int64(0)
	if pageSize > 0 {
		totalPages = (totalCount + int64(pageSize) - 1) / int64(pageSize)
	}

	hasNext := int64(page) < totalPages
	hasPrev := page > 1
	paginatedResponse := &pagination.PaginatedResponse[*models.ListAllRolesResponse]{
		Data: responseRoles,
		Pagination: pagination.PaginationMeta{
			Page:       page,
			PageSize:   pageSize,
			TotalPages: int(totalPages),
			TotalItems: int(totalCount),
			HasNext:    hasNext,
			HasPrev:    hasPrev,
		},
	}

	return paginatedResponse, nil
}
