package services

import (
	"github.com/google/uuid"
	"github.com/vviveksharma/auth/db"
	"github.com/vviveksharma/auth/internal/models"
	"github.com/vviveksharma/auth/internal/repo"
	dbmodels "github.com/vviveksharma/auth/models"
)

type TenantService interface{}

type Tenant struct {
	TenantRepo repo.TenantRepositoryInterface
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
	return nil
}

func (t *Tenant) CreateTenant(req *models.CreateTenantRequest) (resp *models.CreateTenantResponse, err error) {
	_, err = t.TenantRepo.GetUserByEmail(req.Email)
	if err != nil {
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

func (t *Tenant) LoginTenant(req *models.LoginTenantRequest) (resp *models.LoginTenantResponse, err error) {
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
	token := uuid.New()
	return &models.LoginTenantResponse{
		Token: token.String(),
	}, nil
}
