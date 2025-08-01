package repo

import (
	"github.com/google/uuid"
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

type TenantRepositoryInterface interface {
	CreateTenant(tenant *models.DBTenant) error
	GetUserByEmail(email string) (tenantDetails *models.DBTenant, err error)
}

type TenantRepository struct {
	DB *gorm.DB
}

func NewTenantRepository(db *gorm.DB) (TenantRepositoryInterface, error) {
	return &TenantRepository{DB: db}, nil
}

func (t *TenantRepository) CreateTenant(tenant *models.DBTenant) error {
	transaction := t.DB.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	newTenant := transaction.Create(&tenant)
	if newTenant.Error != nil {
		return newTenant.Error
	}
	transaction.Commit()
	return nil
}

func (t *TenantRepository) GetUserByEmail(email string) (tenantDetails *models.DBTenant, err error) {
	transaction := t.DB.Begin()
	if transaction.Error != nil {
		return nil, transaction.Error
	}
	defer transaction.Rollback()
	tenant := transaction.First(&tenantDetails, models.DBUser{
		Email: email,
	})
	if tenant.Error != nil {
		return nil, tenant.Error
	}
	transaction.Commit()
	return tenantDetails, nil
}

func (t *TenantRepository) VerifyTenant(tenantId string) (bool, error) {
	transaction := t.DB.Begin()
	if transaction.Error != nil {
		return false, transaction.Error
	}
	defer transaction.Rollback()
	var tenantDetails models.DBTenant
	findErr := transaction.Model(models.DBTenant{}).Where("id = ?", uuid.MustParse(tenantId)).First(&tenantDetails)
	if findErr.Error != nil {
		return false, findErr.Error
	}
	if tenantDetails.Id == uuid.MustParse(tenantId) {
		return true, nil
	}
	return false, &models.ServiceResponse{
		Code:    400,
		Message: "no tenant found with this id",
	}
}
