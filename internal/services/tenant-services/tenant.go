package tenantservices

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/db"
	"github.com/vviveksharma/auth/internal/models"
	"github.com/vviveksharma/auth/internal/pagination"
	"github.com/vviveksharma/auth/internal/repo"
	"github.com/vviveksharma/auth/internal/utils"
	responsemodels "github.com/vviveksharma/auth/internal/dto/tenant/responses"
	reqmodels "github.com/vviveksharma/auth/internal/dto/tenant/requests"
	dbmodels "github.com/vviveksharma/auth/models"
	smtpservice "github.com/vviveksharma/auth/smtp-service"
)

type ITenantService interface {
	CreateTenant(req *reqmodels.CreateTenantRequest) (resp *responsemodels.CreateTenantResponse, err error)
	LoginTenant(req *reqmodels.LoginTenantRequest, ip string) (resp *responsemodels.LoginTenantResponse, err error)
	ListTokens(ctx context.Context) (resp []*models.ListTokensResponse, err error)
	RevokeToken(ctx context.Context, token string) (resp *responsemodels.RevokeTokenResponse, err error)
	CreateToken(ctx context.Context, req *models.CreateTokenRequest) (*models.CreateTokenResponse, error)
	ResetPassword(ctx context.Context, req *models.ResetTenantPasswordRequest) (*models.ResetPasswordTenantResponse, error)
	SetPassword(ctx context.Context, req *models.SetTenantPasswordRequest) (*models.SetTenantPasswordResponse, error)
	GetTenantDetails(ctx context.Context) (resp *responsemodels.GetTenantDetails, err error)
	DeleteTenant(ctx context.Context) (resp *responsemodels.DeleteTenantResponse, err error)
	ListTokensWithStatus(ctx context.Context, page int, pageSize int, status string) (resp *pagination.PaginatedResponse[*models.GetListTokenWithStatus], err error)
	GetDashboardDetails(ctx context.Context) (resp *responsemodels.DashboardTenantResponse, err error)
}

type Tenant struct {
	TenantRepo      repo.TenantRepositoryInterface
	TokenRepo       repo.TokenRepositoryInterface
	TenantLoginRepo repo.TenantLoginRepositoryInterface
	UserRepo        repo.UserRepositoryInterface
	RoleRepo        repo.RoleRepositoryInterface
	UserLoginRepo   repo.LoginRepositoryInterface
	EmailService    smtpservice.MailServiceInterface
}

func NewTenantService() (ITenantService, error) {
	ser := &Tenant{}
	err := ser.SetupRepo()
	if err != nil {
		return nil, err
	}
	emailService := smtpservice.NewMailService()
	ser.EmailService = emailService
	return ser, nil
}

func (t *Tenant) SetupRepo() error {
	var err error
	tenant, err := repo.NewTenantRepository(db.DB)
	if err != nil {
		return err
	}
	t.TenantRepo = tenant
	token, err := repo.NewTokenRepository(db.DB)
	if err != nil {
		return err
	}
	t.TokenRepo = token
	tenantLogin, err := repo.NewTenantLoginRepository(db.DB)
	if err != nil {
		return err
	}
	t.TenantLoginRepo = tenantLogin
	user, err := repo.NewUserRepository(db.DB)
	if err != nil {
		return err
	}
	t.UserRepo = user
	userLogin, err := repo.NewLoginRepository(db.DB)
	if err != nil {
		return err
	}
	t.UserLoginRepo = userLogin
	repoRole, err := repo.NewRoleRepository(db.DB)
	if err != nil {
		return err
	}
	t.RoleRepo = repoRole
	return nil
}

func (t *Tenant) CreateTenant(req *reqmodels.CreateTenantRequest) (resp *responsemodels.CreateTenantResponse, err error) {
	_, err = t.TenantRepo.GetUserByEmail(req.Email)
	if err != nil && err.Error() != "record not found" {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while checking the tenant details: " + err.Error(),
		}
	}
	// adding the hashing password
	hashedPassword, salt, err := utils.GeneratePasswordHash(req.Password, utils.DefaultParams)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while hashing the password while creating the tenant: " + err.Error(),
		}
	}
	err = t.TenantRepo.CreateTenant(&dbmodels.DBTenant{
		Name:     req.Name,
		Email:    req.Email,
		Campany:  req.Campany,
		Password: hashedPassword,
		Salt:     salt,
	})
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while creating the tenant details: " + err.Error(),
		}
	}
	return &models.CreateTenantResponse{
		Message: "Tenant created successfully for " + req.Name + " for organisation " + req.Campany,
	}, nil
}

func (t *Tenant) LoginTenant(req *reqmodels.LoginTenantRequest, ip string) (resp *responsemodels.LoginTenantResponse, err error) {
	tenantDetails, err := t.TenantRepo.GetUserByEmail(req.Email)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "Tenant with the specified email does not exist.",
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "An error occurred while retrieving tenant details: " + err.Error(),
			}
		}
	}
	tenantLoginDetails, err := t.TenantLoginRepo.GetDetailsByEmail(req.Email)
	if err != nil {
		if err.Error() != "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "An error occurred while retrieving tenant login details: " + err.Error(),
			}
		}
	}
	checkPassword, err := utils.ComparePassword(req.Password, tenantDetails.Password, tenantDetails.Salt, utils.DefaultParams)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "An error occurred while verifying the password: " + err.Error(),
		}
	}
	if !checkPassword {
		return nil, &dbmodels.ServiceResponse{
			Code:    401,
			Message: "The provided password is incorrect. Please try again.",
		}
	}
	var token string
	fmt.Println("the login details: ", tenantLoginDetails)
	if tenantLoginDetails == nil {
		fmt.Println("First time user creating the default and application token")
		tokenName := "logintoken" + utils.GenerateRandomString(5)
		terr := t.TokenRepo.CreateToken(&dbmodels.DBToken{
			TenantId:       tenantDetails.Id,
			Name:           tokenName,
			ExpiresAt:      time.Now().Add(120 * time.Minute),
			IsActive:       true,
			ApplicationKey: false,
			CreatedAt:      time.Now(),
		})
		if terr != nil {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "Failed to create authentication token: " + terr.Error(),
			}
		}
		//Create a default application token
		defaultTokenErr := t.TokenRepo.CreateToken(&dbmodels.DBToken{
			TenantId:       tenantDetails.Id,
			Name:           "defaultToken",
			ExpiresAt:      time.Now().Add(120 * time.Minute),
			IsActive:       true,
			ApplicationKey: true,
			CreatedAt:      time.Now(),
		})
		if defaultTokenErr != nil {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "Failed to create authentication token: " + defaultTokenErr.Error(),
			}
		}
		tlerr := t.TenantLoginRepo.Create(&dbmodels.DBTenantLogin{
			TenantId:  tenantDetails.Id,
			ExpiresAt: time.Now().Add(24 * time.Hour),
			IsActive:  true,
			IPAddress: ip,
			Email:     req.Email,
		})
		if tlerr != nil {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "Failed to create tenant login record: " + tlerr.Error(),
			}
		}
		// giving the login token as respose to be used further
		tokenDetails, err := t.TokenRepo.GetTokenDetails(dbmodels.DBToken{Name: tokenName})
		if err != nil {
			if err.Error() == "record not found" {
				return nil, &dbmodels.ServiceResponse{
					Code:    404,
					Message: "token with this doesn't exist",
				}
			} else {
				return nil, &dbmodels.ServiceResponse{
					Code:    500,
					Message: "error while fetching the tokendetails: " + err.Error(),
				}
			}
		}
		token = tokenDetails.Id.String()
	} else {
		fmt.Println("Updating the login token")
		newToken, err := t.TokenRepo.UpdateLoginToken(tenantDetails.Id)
		if err != nil {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "Failed to update the login token: " + err.Error(),
			}
		}
		token = newToken.String()
	}
	return &models.LoginTenantResponse{
		Token: token,
	}, nil
}

func (t *Tenant) ListTokens(ctx context.Context) (resp []*models.ListTokensResponse, err error) {
	tenantId := ctx.Value("tenant_id").(string)
	dbTokens, err := t.TokenRepo.ListTokens(uuid.MustParse(tenantId))
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while listing tokens for the given tenant :" + err.Error(),
		}
	}
	for _, token := range dbTokens {
		if token.ApplicationKey {
			tokenDetails := models.ListTokensResponse{
				CreateAt:  token.CreatedAt,
				ExpiresAt: token.ExpiresAt,
				Name:      token.Name,
				TokenId:   token.Id,
				Status:    token.IsActive,
			}
			resp = append(resp, &tokenDetails)
		}
	}
	return resp, nil
}

func (t *Tenant) RevokeToken(ctx context.Context, tokenId string) (resp *responsemodels.RevokeTokenResponse, err error) {
	tokenDetails, err := t.TokenRepo.GetTokenDetails(dbmodels.DBToken{
		Id: uuid.MustParse(tokenId),
	})
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "Token not found",
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "Failed to retrieve token details",
			}
		}
	}
	if !tokenDetails.IsActive {
		return nil, &dbmodels.ServiceResponse{
			Code:    409,
			Message: "Token is already inactive",
		}
	}
	err = t.TokenRepo.RevokeToken(tokenId)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "Failed to revoke token",
		}
	}
	return &models.RevokeTokenResponse{
		Message: "Token revoked successfully",
	}, nil
}

func (t *Tenant) CreateToken(ctx context.Context, req *models.CreateTokenRequest) (*models.CreateTokenResponse, error) {
	tenantId := ctx.Value("tenant_id").(string)
	tokendata, err := t.TokenRepo.GetTokenDetails(dbmodels.DBToken{Name: req.Name})
	if err != nil {
		if err.Error() != "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "error while fetching the token details :" + err.Error(),
			}
		}
	}
	log.Println("the token data: ", tokendata)
	if tokendata != nil && tokendata.Name == req.Name {
		return nil, &dbmodels.ServiceResponse{
			Code:    409,
			Message: "record already exists please try with another name",
		}
	}
	parsedExpiry, parseErr := time.Parse("2006-01-02", req.ExpiryAt)
	if parseErr != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    422,
			Message: "Invalid expiry date format. Please use YYYY-MM-DD.",
		}
	}
	err = t.TokenRepo.CreateToken(&dbmodels.DBToken{
		TenantId:       uuid.MustParse(tenantId),
		Name:           req.Name,
		CreatedAt:      time.Now(),
		ExpiresAt:      parsedExpiry,
		ApplicationKey: true,
		IsActive:       true,
	})
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "erorr while creating a token: " + err.Error(),
		}
	}
	return &models.CreateTokenResponse{
		Message: "Token created successfully",
	}, nil
}

func (t *Tenant) ResetPassword(ctx context.Context, req *models.ResetTenantPasswordRequest) (*models.ResetPasswordTenantResponse, error) {
	tenantDetails, err := t.TenantRepo.GetUserByEmail(req.Email)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "tenant with this email doesnot exist",
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "error while fetching the tenant details: " + err.Error(),
			}
		}
	}
	log.Println(tenantDetails)
	emailErr := t.EmailService.SendEmailWithTemplate(req.Email, "Reset Your Password", "password_reset", map[string]interface{}{
		"name":       tenantDetails.Name,
		"reset_link": "http://localhost:8080/tenant/setpassword",
	})
	if emailErr != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while sending email to the user: " + emailErr.Error(),
		}
	}
	return &models.ResetPasswordTenantResponse{
		Message: "Email sent successfully to email " + req.Email,
	}, nil
}

func (t *Tenant) SetPassword(ctx context.Context, req *models.SetTenantPasswordRequest) (*models.SetTenantPasswordResponse, error) {
	tenantDetails, err := t.TenantRepo.GetUserByEmail(req.Email)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "Tenant with the specified email does not exist.",
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "An error occurred while retrieving tenant details: " + err.Error(),
			}
		}
	}
	hashedPassword, err := utils.GeneratePassword(req.NewPassword, utils.DefaultParams, tenantDetails.Salt)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "An error occurred while generating the new password hash: " + err.Error(),
		}
	}
	log.Println("the stored the hashed password: ", tenantDetails.Password)
	log.Println("the generated one: ", hashedPassword)
	// Ensure the new password is different from the old password
	if hashedPassword == tenantDetails.Password {
		return nil, &dbmodels.ServiceResponse{
			Code:    409,
			Message: "The new password cannot be the same as the current password.",
		}
	}
	// Update the password in the database
	err = t.TenantRepo.UpdateTenatDetailsPassword(tenantDetails.Id.String(), hashedPassword)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "An error occurred while updating the password: " + err.Error(),
		}
	}
	return &models.SetTenantPasswordResponse{
		Message: "Password has been updated successfully.",
	}, nil
}

func (t *Tenant) GetTenantDetails(ctx context.Context) (resp *responsemodels.GetTenantDetails, err error) {
	tenantId := ctx.Value("tenant_id").(string)
	tenantDetails, err := t.TenantRepo.GetTenantDetails(&dbmodels.DBTenant{
		Id: uuid.MustParse(tenantId),
	})
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "tenant with this email doesnot exist",
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "error while fetching the tenant details: " + err.Error(),
			}
		}
	}
	return &models.GetTenantDetails{
		Name:         tenantDetails.Name,
		Email:        tenantDetails.Email,
		Organisation: tenantDetails.Campany,
	}, nil
}

func (t *Tenant) DeleteTenant(ctx context.Context) (resp *responsemodels.DeleteTenantResponse, err error) {
	// deleting the tenant
	tenantId := ctx.Value("tenant_id").(string)
	log.Printf("Attempting to delete tenant with ID: %s", tenantId)

	err = t.TenantRepo.DeleteTenant(uuid.MustParse(tenantId))
	if err != nil {
		if err.Error() == "record not found" {
			log.Printf("Delete failed: tenant not found for ID: %s", tenantId)
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "tenant with this ID does not exist",
			}
		} else {
			log.Printf("Delete failed: error deleting tenant ID %s: %v", tenantId, err)
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "error while deleting the tenant: " + err.Error(),
			}
		}
	}

	log.Printf("Successfully deleted tenant with ID: %s", tenantId)
	return &models.DeleteTenantResponse{
		Message: "tenant deleted successfully",
	}, nil
}

func (t *Tenant) GetDashboardDetails(ctx context.Context) (resp *responsemodels.DashboardTenantResponse, err error) {
	tenantId := ctx.Value("tenant_id").(string)
	userDetails, _, err := t.UserRepo.ListUsersPaginated(1, 1000, uuid.MustParse(tenantId), "enabled")
	if err != nil {
		if err.Error() == "record not found" {
			log.Printf("Delete failed: tenant not found for ID: %s", tenantId)
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "tenant with this ID does not exist",
			}
		} else {
			log.Printf("Delete failed: error deleting tenant ID %s: %v", tenantId, err)
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "error while fetching the user of the tenant: " + err.Error(),
			}
		}
	}

	roleDetails, err := t.RoleRepo.GetRolesByTenant(uuid.MustParse(tenantId), "custom")
	if err != nil {
		if err.Error() == "record not found" {
			log.Printf("Delete failed: tenant not found for ID: %s", tenantId)
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "role details tenant with this ID does not exist",
			}
		} else {
			log.Printf("Delete failed: error deleting tenant ID %s: %v", tenantId, err)
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "error while fetching the role of the tenant: " + err.Error(),
			}
		}
	}
	tokenDetails, err := t.TokenRepo.ListTokens(uuid.MustParse(tenantId))
	if err != nil {
		if err.Error() == "record not found" {
			log.Printf("Delete failed: tenant not found for ID: %s", tenantId)
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "role details tenant with this ID does not exist",
			}
		} else {
			log.Printf("Delete failed: error deleting tenant ID %s: %v", tenantId, err)
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "error while fetching the role of the tenant: " + err.Error(),
			}
		}
	}
	var tokensize int
	for _, token := range tokenDetails {
		if token.ApplicationKey {
			tokensize += 1
		}
	}
	return &models.DashboardTenantResponse{
		UsersCount: len(userDetails),
		RoleCount:  len(roleDetails),
		TokenCount: tokensize,
	}, nil
}

func (t *Tenant) ListTokensWithStatus(ctx context.Context, page int, pageSize int, status string) (resp *pagination.PaginatedResponse[*models.GetListTokenWithStatus], err error) {
	tenantId := ctx.Value("tenant_id").(string)
	tokens, totalCount, err := t.TokenRepo.ListTokensPaginated(uuid.MustParse(tenantId), page, pageSize, status)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while listing tokens: " + err.Error(),
		}
	}
	var responseTokens []*models.GetListTokenWithStatus
	for _, t := range tokens {
		token := &models.GetListTokenWithStatus{
			CreateAt:  t.CreatedAt,
			ExpiresAt: t.ExpiresAt,
			TokenId:   t.Id,
			Status:    t.IsActive,
			Name:      t.Name,
		}
		responseTokens = append(responseTokens, token)
	}
	totalPages := int64(0)
	if pageSize > 0 {
		totalPages = (totalCount + int64(pageSize) - 1) / int64(pageSize)
	}

	hasNext := int64(page) < totalPages
	hasPrev := page > 1
	paginatedResponse := &pagination.PaginatedResponse[*models.GetListTokenWithStatus]{
		Data: responseTokens,
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
