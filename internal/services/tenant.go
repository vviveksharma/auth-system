package services

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/db"
	"github.com/vviveksharma/auth/internal/models"
	"github.com/vviveksharma/auth/internal/repo"
	"github.com/vviveksharma/auth/internal/utils"
	dbmodels "github.com/vviveksharma/auth/models"
)

type TenantService interface {
	CreateTenant(req *models.CreateTenantRequest) (resp *models.CreateTenantResponse, err error)
	LoginTenant(req *models.LoginTenantRequest, ip string) (resp *models.LoginTenantResponse, err error)
	ListTokens(ctx context.Context, logintoken string) (resp []*models.ListTokensResponse, err error)
	RevokeToken(ctx context.Context, token string) (resp *models.RevokeTokenResponse, err error)
	CreateToken(ctx context.Context, req *models.CreateTokenRequest) (*models.CreateTokenResponse, error)
}

type Tenant struct {
	TenantRepo      repo.TenantRepositoryInterface
	TokenRepo       repo.TokenRepositoryInterface
	TenantLoginRepo repo.TenantLoginRepositoryInterface
}

func NewTenantService() (TenantService, error) {
	ser := &Tenant{}
	err := ser.SetupRepo()
	if err != nil {
		return nil, err
	}
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
	err = t.TenantRepo.CreateTenant(&dbmodels.DBTenant{
		Name:     req.Name,
		Email:    req.Email,
		Campany:  req.Campany,
		Password: req.Password,
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
				Message: "user with name doesnot exist",
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "error while creating the userdetails: " + err.Error(),
			}
		}
	}
	if req.Password != tenantDetails.Password {
		return nil, &dbmodels.ServiceResponse{
			Code:    401,
			Message: "invalid email or password",
		}
	}
	token := uuid.New().String()
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
			Message: "error while creating a token: " + terr.Error(),
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
			Message: "error while creating tenant login: " + tlerr.Error(),
		}
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

func (t *Tenant) RevokeToken(ctx context.Context, token string) (resp *models.RevokeTokenResponse, err error) {
	err = t.TokenRepo.RevokeToken(token)
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
	tokendata,err := t.TokenRepo.GetTokenDetailsByName(req.Name)
	if err != nil {
		if err.Error() != "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code: 500,
				Message: "error while fetching the token details :" + err.Error(),
			}
		}
	}
	log.Println("the token data: ",tokendata)
	if tokendata != nil &&  tokendata.Name == req.Name {
		return nil, &dbmodels.ServiceResponse{
			Code: 423,
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
