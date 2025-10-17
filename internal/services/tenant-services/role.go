package tenantservices

import (
	"context"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/db"
	tenantmodels "github.com/vviveksharma/auth/internal/models/tenantModels"
	"github.com/vviveksharma/auth/internal/pagination"
	tenantrepo "github.com/vviveksharma/auth/internal/repo/tenantRepo"
	dbmodels "github.com/vviveksharma/auth/models"
)

type ITenantRoleService interface{}

type TenantRoleService struct {
	RoleRepo tenantrepo.TenantRoleRepositoryInterface
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
	repoRole, err := tenantrepo.NewTenantRoleRepository(db.DB)
	if err != nil {
		return err
	}
	tr.RoleRepo = repoRole
	return nil
}

func (tr *TenantRoleService) ListRoles(ctx context.Context, page int, page_size int, roleType string, status string) (resp *pagination.PaginatedResponse[*tenantmodels.TenantListRoleResponseBody], err error) {
	tenantId := ctx.Value("tenant_id").(string)
	roles, totalCount, err := tr.RoleRepo.ListRoles(uuid.MustParse(tenantId), page, page_size, status, roleType)
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
