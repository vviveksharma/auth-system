package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/db"
	"github.com/vviveksharma/auth/internal/models"
	"github.com/vviveksharma/auth/internal/repo"
	"github.com/vviveksharma/auth/internal/utils"
	dbmodels "github.com/vviveksharma/auth/models"
	smtpservice "github.com/vviveksharma/auth/smtp-service"
)

type TenantService interface {
	CreateTenant(req *models.CreateTenantRequest) (resp *models.CreateTenantResponse, err error)
	LoginTenant(req *models.LoginTenantRequest, ip string) (resp *models.LoginTenantResponse, err error)
	ListTokens(ctx context.Context, logintoken string) (resp []*models.ListTokensResponse, err error)
	RevokeToken(ctx context.Context, token string) (resp *models.RevokeTokenResponse, err error)
	CreateToken(ctx context.Context, req *models.CreateTokenRequest) (*models.CreateTokenResponse, error)
	ResetPassword(ctx context.Context, req *models.ResetTenantPasswordRequest) (*models.ResetPasswordTenantResponse, error)
	SetPassword(ctx context.Context, req *models.SetTenantPasswordRequest) (*models.SetTenantPasswordResponse, error)
	ListUsers(ctx context.Context) (resp []*models.ListUserTenant, err error)
}

type Tenant struct {
	TenantRepo      repo.TenantRepositoryInterface
	TokenRepo       repo.TokenRepositoryInterface
	TenantLoginRepo repo.TenantLoginRepositoryInterface
	UserRepo        repo.UserRepositoryInterface
	UserLoginRepo   repo.LoginRepositoryInterface
	EmailService    smtpservice.MailServiceInterface
}

func NewTenantService() (TenantService, error) {
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
	return nil
}

func (t *Tenant) CreateTenant(req *models.CreateTenantRequest) (resp *models.CreateTenantResponse, err error) {
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

func (t *Tenant) LoginTenant(req *models.LoginTenantRequest, ip string) (resp *models.LoginTenantResponse, err error) {
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
	token := uuid.New().String()
	if tenantLoginDetails == nil {
		terr := t.TokenRepo.CreateToken(&dbmodels.DBToken{
			TenantId:       tenantDetails.Id,
			Name:           "logintoken" + utils.GenerateRandomString(5),
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
		})
		if tlerr != nil {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "Failed to create tenant login record: " + tlerr.Error(),
			}
		}
	} else {
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

func (t *Tenant) ListTokens(ctx context.Context, logintoken string) (resp []*models.ListTokensResponse, err error) {
	tenantId, err := t.TokenRepo.GetTenantUsingToken(logintoken)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while fetching the tenant details: " + err.Error(),
		}
	}
	dbTokens, err := t.TokenRepo.ListTokens(*tenantId)
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
			}
			resp = append(resp, &tokenDetails)
		}
	}
	return resp, nil
}

func (t *Tenant) RevokeToken(ctx context.Context, tokenId string) (resp *models.RevokeTokenResponse, err error) {
	err = t.TokenRepo.RevokeToken(tokenId)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while revoking the token: " + err.Error(),
		}
	}
	return &models.RevokeTokenResponse{
		Message: "Token revoked successfully.",
	}, nil
}

func (t *Tenant) CreateToken(ctx context.Context, req *models.CreateTokenRequest) (*models.CreateTokenResponse, error) {
	tenantId := ctx.Value("tenant_id").(string)
	tokendata, err := t.TokenRepo.GetTokenDetailsByName(req.Name)
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
			Code:    423,
			Message: "record already exists please try with another name",
		}
	}
	parsedExpiry, parseErr := time.Parse("2006-01-02", req.ExpiryAt)
	if parseErr != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    400,
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
			Code:    423,
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

func (t *Tenant) ListUsers(ctx context.Context) (resp []*models.ListUserTenant, err error) {
	tenantId := ctx.Value("tenant_id").(string)
	fmt.Println(tenantId)
	userDetails, err := t.UserRepo.ListUsers(uuid.MustParse(tenantId))
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while fetching user data for the particular tenant: " + err.Error(),
		}
	}
	loginDetails, err := t.UserLoginRepo.GetUsers(uuid.MustParse(tenantId))
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while fetching login data for the particular tenant: " + err.Error(),
		}
	}
	var user models.ListUserTenant
	for _, users := range userDetails {
		for _, login := range loginDetails {
			user.Email = users.Name
			user.Name = users.Name
			user.CreatedAt = users.CreatedAt.String()
			user.Role = login.RoleName
			user.LogginStatus = login.Revoked
			resp = append(resp, &user)
		}
	}
	return resp, nil
}
