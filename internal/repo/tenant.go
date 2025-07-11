package repo

import (
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

type TenantRepositoryInterface interface{
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