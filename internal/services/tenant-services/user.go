package tenantservices

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/db"
	tenantmodels "github.com/vviveksharma/auth/internal/dto/tenant/responses"
	"github.com/vviveksharma/auth/internal/pagination"
	"github.com/vviveksharma/auth/internal/repo"
	tenantrepo "github.com/vviveksharma/auth/internal/repo/tenantRepo"
	dbmodels "github.com/vviveksharma/auth/models"
)

type ITenantUserService interface {
	TenantListUsers(ctx context.Context, page int, pageSize int, stautus string) (resp *pagination.PaginatedResponse[*tenantmodels.TenantListUserResponseBody], err error)
	TenantEnableUser(ctx context.Context, userId uuid.UUID) (resp *tenantmodels.TenantEnableUserResponseBody, err error)
	TenantDisableUser(ctx context.Context, userId uuid.UUID) (resp *tenantmodels.TenantEnableUserResponseBody, err error)
	TenantDeleteUser(ctx context.Context, userId uuid.UUID) (resp *tenantmodels.TenantDeleteUserResponseBody, err error)
}

type TenantUserService struct {
	TenantUserRepo tenantrepo.TenantUserRepositoryInterface
	UserRepo       repo.UserRepositoryInterface
	SharedRepo     repo.SharedRepoInterface
}

func NewTenantUserService() (ITenantUserService, error) {
	ser := &TenantUserService{}
	err := ser.SetupRepo()
	if err != nil {
		return nil, err
	}
	return ser, nil
}

func (tu *TenantUserService) SetupRepo() error {
	tenantUserRepo, err := tenantrepo.NewTenantUserRepository(db.DB)
	if err != nil {
		return err
	}
	tu.TenantUserRepo = tenantUserRepo
	userRepo, err := repo.NewUserRepository(db.DB)
	if err != nil {
		return err
	}
	tu.UserRepo = userRepo
	sharedRepo, err := repo.NewSharedRepository(db.DB)
	if err != nil {
		return err
	}
	tu.SharedRepo = sharedRepo
	return nil
}

func (tu *TenantUserService) TenantListUsers(ctx context.Context, page int, pageSize int, stautus string) (resp *pagination.PaginatedResponse[*tenantmodels.TenantListUserResponseBody], err error) {
	tenantId := ctx.Value("tenant_id").(string)
	log.Printf("🔍 [TENANT SERVICE] TenantListUsers called with: page=%d, pageSize=%d, status=%s, tenant_id=%s",
		page, pageSize, stautus, tenantId)

	userDetails, count, err := tu.TenantUserRepo.ListUsers(page, pageSize, uuid.MustParse(tenantId), stautus)
	if err != nil {
		log.Printf("❌ [TENANT SERVICE] Error from repo: %v", err)
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while rendering the tenant roles based filtered response: " + err.Error(),
		}
	}

	log.Printf("✅ [TENANT SERVICE] Received from repo: %d users, totalCount=%d", len(userDetails), count)

	var response []*tenantmodels.TenantListUserResponseBody
	for i, users := range userDetails {
		log.Printf("   [TENANT SERVICE] Processing user %d: ID=%s, Email=%s, Name=%s, Status=%v, Roles=%v",
			i+1, users.Id.String(), users.Email, users.Name, users.Status, users.Roles)

		user := &tenantmodels.TenantListUserResponseBody{
			Id:        users.Id,
			Name:      users.Name,
			Status:    users.Status,
			CreatedAt: users.CreatedAt,
			Roles:     users.Roles,
			Email:     users.Email,
		}
		response = append(response, user)
	}

	log.Printf("🔍 [TENANT SERVICE] Built response array with %d users", len(response))

	totalPages := int64(0)
	if pageSize > 0 {
		totalPages = (count + int64(pageSize) - 1) / int64(pageSize)
	}
	hasNext := int64(page) < totalPages
	hasPrev := page > 1

	log.Printf("🔍 [TENANT SERVICE] Pagination: totalPages=%d, hasNext=%v, hasPrev=%v", totalPages, hasNext, hasPrev)

	paginatedResponse := &pagination.PaginatedResponse[*tenantmodels.TenantListUserResponseBody]{
		Data: response,
		Pagination: pagination.PaginationMeta{
			Page:       page,
			PageSize:   pageSize,
			TotalPages: int(totalPages),
			TotalItems: int(count),
			HasNext:    hasNext,
			HasPrev:    hasPrev,
		},
	}

	log.Printf("✅ [TENANT SERVICE] Returning paginated response with %d users", len(paginatedResponse.Data))
	return paginatedResponse, nil
}

func (tu *TenantUserService) TenantEnableUser(ctx context.Context, userId uuid.UUID) (resp *tenantmodels.TenantEnableUserResponseBody, err error) {
	tenantId := ctx.Value("tenant_id").(string)
	userDetails, err := tu.UserRepo.GetUserDetails(dbmodels.DBUser{
		TenantId: uuid.MustParse(tenantId),
		Id:       userId,
	})
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "User with this id doesn't exist",
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "error while finding the userDetails",
			}
		}
	}
	if userDetails.Status {
		return nil, &dbmodels.ServiceResponse{
			Code:    409,
			Message: "User is already enabled",
		}
	}
	err = tu.UserRepo.ChangeStatus(true, userId)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while updating the user status",
		}
	}
	return &tenantmodels.TenantEnableUserResponseBody{
		Message: fmt.Sprintf("user with userId %s enabled successfully", userId.String()),
	}, nil
}

func (tu *TenantUserService) TenantDisableUser(ctx context.Context, userId uuid.UUID) (resp *tenantmodels.TenantEnableUserResponseBody, err error) {
	tenantId := ctx.Value("tenant_id").(string)
	userDetails, err := tu.UserRepo.GetUserDetails(dbmodels.DBUser{
		TenantId: uuid.MustParse(tenantId),
		Id:       userId,
	})
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "User with this id doesn't exist",
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "error while finding the userDetails",
			}
		}
	}
	if !userDetails.Status {
		return nil, &dbmodels.ServiceResponse{
			Code:    409,
			Message: "User is already disbaled",
		}
	}
	err = tu.UserRepo.ChangeStatus(true, userId)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while updating the user status",
		}
	}
	return &tenantmodels.TenantEnableUserResponseBody{
		Message: fmt.Sprintf("user with userId %s disabled successfully", userId.String()),
	}, nil
}

func (tu *TenantUserService) TenantDeleteUser(ctx context.Context, userId uuid.UUID) (resp *tenantmodels.TenantDeleteUserResponseBody, err error) {
	tenantId := ctx.Value("tenant_id").(string)
	userDetails, err := tu.UserRepo.GetUserDetails(dbmodels.DBUser{
		TenantId: uuid.MustParse(tenantId),
		Id:       userId,
	})
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "User with this id doesn't exist",
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "error while finding the userDetails",
			}
		}
	}
	fmt.Println("userDetails: ", userDetails)
	err = tu.SharedRepo.DeleteUser(userId, uuid.MustParse(tenantId))
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while deleting the user from the database: " + err.Error(),
		}
	}
	return &tenantmodels.TenantDeleteUserResponseBody{
		Message: fmt.Sprintf("successfully deleted user with the user-id %s", userId),
	}, nil
}
