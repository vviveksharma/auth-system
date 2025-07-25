package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/db"
	"github.com/vviveksharma/auth/internal/models"
	"github.com/vviveksharma/auth/internal/repo"
	dbmodels "github.com/vviveksharma/auth/models"
)

type TenantService interface {
	CreateTenant(req *models.CreateTenantRequest) (resp *models.CreateTenantResponse, err error)
	LoginTenant(req *models.LoginTenantRequest, ip string) (resp *models.LoginTenantResponse, err error)
	ListTokens(ctx context.Context, token string) ([]*string, error)
	RevokeToken(ctx context.Context, token string) (resp *models.RevokeTokenResponse, err error)
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
		TenantId:  tenantDetails.Id,
		Token:     token,
		ExpiresAt: time.Now().Add(120 * time.Minute),
		IsActive:  true,
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

func (t *Tenant) ListTokens(ctx context.Context, token string) ([]*string, error) {
	tenantId, err := t.TokenRepo.GetTenantUsingToken(token)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{}
	}
	dbTokens, err := t.TokenRepo.ListTokens(*tenantId)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{}
	}
	var tokens []*string
	for _, dbToken := range dbTokens {
		if !dbToken.IsActive {
			continue
		}
		tokens = append(tokens, &dbToken.Token)
	}
	return tokens, nil
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
