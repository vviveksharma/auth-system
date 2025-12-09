package tenantservices

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/db"
	"github.com/vviveksharma/auth/internal/models"
	tenantmodels "github.com/vviveksharma/auth/internal/models/tenantModels"
	"github.com/vviveksharma/auth/internal/pagination"
	"github.com/vviveksharma/auth/internal/repo"
	tenantrepo "github.com/vviveksharma/auth/internal/repo/tenantRepo"
	dbmodels "github.com/vviveksharma/auth/models"
)

type ITenantRoleService interface {
	TenantListRoles(ctx context.Context, page int, page_size int, roleType string, status string) (resp *pagination.PaginatedResponse[*tenantmodels.TenantListRoleResponseBody], err error)
	TenantGetRolePermissions(ctx context.Context, roleId uuid.UUID) (resp *tenantmodels.TenantGetPermissionsResponseBody, err error)
	TenantDisableRole(ctx context.Context, roleId uuid.UUID) (resp *tenantmodels.TenantDisableRoleResponsBody, err error)
	TenantEnableRole(ctx context.Context, roleId uuid.UUID) (resp *tenantmodels.TenantEnableRoleResponseBody, err error)
	TenantEditRolePermissions(ctx context.Context, roleId uuid.UUID, req *tenantmodels.TeanantEditPermissionRequestBody) (resp *tenantmodels.TenantEditPermissionRoleResponseBody, err error)
	TenantAddRole(ctx context.Context, req *tenantmodels.TenantAddRoleRequestBody) (resp *tenantmodels.TenantAddRoleResponseBody, err error)
	TenantDeleteRole(ctx context.Context, roleId uuid.UUID) (resp *tenantmodels.TenantDeleteRoleResponseBody, err error)
}

type TenantRoleService struct {
	TenantRoleRepo tenantrepo.TenantRoleRepositoryInterface
	RoleRepo       repo.RoleRepositoryInterface
	RouteRoleRepo  repo.RouteRoleRepositoryInterface
	SharedRepo     repo.SharedRepoInterface
}

func NewTenantRoleService() (ITenantRoleService, error) {
	ser := &TenantRoleService{}
	err := ser.SetupRepo()
	if err != nil {
		return nil, err
	}
	return ser, nil
}

func (tr *TenantRoleService) SetupRepo() error {
	var err error
	tenantRepoRole, err := tenantrepo.NewTenantRoleRepository(db.DB)
	if err != nil {
		return err
	}
	tr.TenantRoleRepo = tenantRepoRole
	roleRepo, err := repo.NewRoleRepository(db.DB)
	if err != nil {
		return err
	}
	tr.RoleRepo = roleRepo
	sharedRepo, err := repo.NewSharedRepository(db.DB)
	if err != nil {
		return err
	}
	tr.SharedRepo = sharedRepo
	return nil
}

func (tr *TenantRoleService) TenantListRoles(ctx context.Context, page int, page_size int, roleType string, status string) (resp *pagination.PaginatedResponse[*tenantmodels.TenantListRoleResponseBody], err error) {
	tenantId := ctx.Value("tenant_id").(string)
	roles, totalCount, err := tr.TenantRoleRepo.ListRoles(uuid.MustParse(tenantId), page, page_size, status, roleType)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while rendering the tenant roles based filtered response: " + err.Error(),
		}
	}
	var response []*tenantmodels.TenantListRoleResponseBody
	for _, role := range roles {
		roleResponse := &tenantmodels.TenantListRoleResponseBody{
			Id:          role.RoleId,
			Name:        role.Role,
			DisplayName: role.DisplayName,
			RoleType:    role.RoleType,
			Status:      role.Status,
		}
		response = append(response, roleResponse)
	}
	totalPages := int64(0)
	if page_size > 0 {
		totalPages = (totalCount + int64(page_size) - 1) / int64(page_size)
	}
	hasNext := int64(page) < totalPages
	hasPrev := page > 1
	paginatedResponse := &pagination.PaginatedResponse[*tenantmodels.TenantListRoleResponseBody]{
		Data: response,
		Pagination: pagination.PaginationMeta{
			Page:       page,
			PageSize:   page_size,
			TotalPages: int(totalPages),
			TotalItems: int(totalCount),
			HasNext:    hasNext,
			HasPrev:    hasPrev,
		},
	}

	return paginatedResponse, nil
}

func (tr *TenantRoleService) TenantGetRolePermissions(ctx context.Context, roleId uuid.UUID) (resp *tenantmodels.TenantGetPermissionsResponseBody, err error) {
	tenantId := ctx.Value("tenant_id").(string)
	isSystemRole := models.IsSystemRole(roleId)
	if isSystemRole {
		tenantId = dbmodels.GetSystemTenantId()
	}
	permissions, err := tr.TenantRoleRepo.GetPermissions(uuid.MustParse(tenantId), roleId)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    404,
			Message: "error while finding the role permissions:" + err.Error(),
		}
	}
	return &tenantmodels.TenantGetPermissionsResponseBody{
		Id:          roleId,
		RoleInfo:    permissions.RoleInfo,
		Permissions: permissions.Permissions,
	}, nil
}

func (tr *TenantRoleService) TenantDisableRole(ctx context.Context, roleId uuid.UUID) (resp *tenantmodels.TenantDisableRoleResponsBody, err error) {
	tenantId := ctx.Value("tenant_id").(string)
	roleDetails, err := tr.RoleRepo.GetRolesDetails(&dbmodels.DBRoles{
		TenantId: uuid.MustParse(tenantId),
		RoleId:   roleId,
	})
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "No role with this role-id exists",
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "error while finding the role:" + err.Error(),
			}
		}
	}
	if !roleDetails.Status {
		return nil, &dbmodels.ServiceResponse{
			Code:    409,
			Message: "the role you are trying to disable is already disabled",
		}
	}
	err = tr.RoleRepo.ChangeStatus(false, roleId)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "Failed to enable the role due to a database error: " + err.Error() + ". Please try again or contact support if the issue persists.",
		}
	}
	return &tenantmodels.TenantDisableRoleResponsBody{
		Message: fmt.Sprintf("the role with the role-id %v is updated successfully", roleId),
	}, nil
}

func (tr *TenantRoleService) TenantEnableRole(ctx context.Context, roleId uuid.UUID) (resp *tenantmodels.TenantEnableRoleResponseBody, err error) {
	tenantId := ctx.Value("tenant_id").(string)
	roleDetails, err := tr.RoleRepo.GetRolesDetails(&dbmodels.DBRoles{
		TenantId: uuid.MustParse(tenantId),
		RoleId:   roleId,
	})
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "No role with this role-id exists",
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "error while finding the role:" + err.Error(),
			}
		}
	}
	if roleDetails.Status {
		return nil, &dbmodels.ServiceResponse{
			Code:    409,
			Message: "the role you are trying to enable is already enabled",
		}
	}
	err = tr.RoleRepo.ChangeStatus(true, roleId)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "Failed to enable the role due to a database error: " + err.Error() + ". Please try again or contact support if the issue persists.",
		}
	}
	return &tenantmodels.TenantEnableRoleResponseBody{
		Message: fmt.Sprintf("the role with the role-id %v is updated successfully", roleId),
	}, nil
}

func (tr *TenantRoleService) TenantEditRolePermissions(ctx context.Context, roleId uuid.UUID, req *tenantmodels.TeanantEditPermissionRequestBody) (resp *tenantmodels.TenantEditPermissionRoleResponseBody, err error) {
	tenantId := ctx.Value("tenant_id").(string)
	roleDetails, err := tr.RoleRepo.GetRolesDetails(&dbmodels.DBRoles{
		TenantId: uuid.MustParse(tenantId),
		RoleId:   roleId,
	})
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "No role with this role-id exists",
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "error while finding the role:" + err.Error(),
			}
		}
	}
	if req.UpdateRoleDetails {
		err = tr.RoleRepo.UpdateRoleDetails(dbmodels.DBRoles{
			Id:          roleDetails.Id,
			Role:        req.RoleInfo.Name,
			DisplayName: req.RoleInfo.DisplayName,
			Description: req.RoleInfo.Description,
			RoleId:      roleDetails.RoleId,
			TenantId:    uuid.MustParse(tenantId),
			RoleType:    "custom",
			Status:      true,
			CreatedAt:   roleDetails.CreatedAt,
			UpdatedAt:   time.Now(),
		}, uuid.MustParse(tenantId))
		if err != nil {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "Failed to update the role due to a database error: " + err.Error(),
			}
		}
	}
	if req.UpdateRolePermissions {
		err := tr.TenantRoleRepo.UpdateRolePermissions(uuid.MustParse(tenantId), roleId, req.Permissions)
		if err != nil {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("Failed to update permissions for role %s due to a database error: %s. Please try again or contact support if the issue persists.", roleId, err.Error()),
			}
		}
	}
	return &tenantmodels.TenantEditPermissionRoleResponseBody{
		Message: fmt.Sprintf("the role with the role-id %v is updated successfully", roleId),
	}, nil
}

func (tr *TenantRoleService) TenantAddRole(ctx context.Context, req *tenantmodels.TenantAddRoleRequestBody) (resp *tenantmodels.TenantAddRoleResponseBody, err error) {
	tenantId := ctx.Value("tenant_id").(string)
	roleDetails, err := tr.RoleRepo.GetRolesDetails(&dbmodels.DBRoles{
		DisplayName: req.DisplayName,
		Role:        req.Name,
		TenantId:    uuid.MustParse(tenantId),
	})
	if err != nil && err.Error() != "record not found" {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while fetching the role details from the database: " + err.Error(),
		}
	}
	if roleDetails != nil {
		if roleDetails.Role == req.Name || roleDetails.DisplayName == req.DisplayName {
			return nil, &dbmodels.ServiceResponse{
				Code:    409,
				Message: "error while creating the role with this name they already exist",
			}
		}
	}
	err = tr.SharedRepo.CreateCustomRole(&models.CreateCustomRole{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Permissions: req.Permissions,
	}, uuid.MustParse(tenantId))
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while creating a role" + err.Error(),
		}
	}
	return &tenantmodels.TenantAddRoleResponseBody{
		Message: fmt.Sprintf("Role with the name %s added successfully", req.Name),
	}, nil
}

func (tr *TenantRoleService) TenantDeleteRole(ctx context.Context, roleId uuid.UUID) (resp *tenantmodels.TenantDeleteRoleResponseBody, err error) {
	tenantId := ctx.Value("tenant_id").(string)
	err = tr.SharedRepo.DeleteCustomRole(roleId, uuid.MustParse(tenantId))
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while deleting the role: " + err.Error(),
		}
	}
	return &tenantmodels.TenantDeleteRoleResponseBody{
		Message: fmt.Sprintf("role with the roleId %s deleted successfully", roleId),
	}, nil
}
