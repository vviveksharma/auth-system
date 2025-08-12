package repo

import (
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

type TenantLoginRepositoryInterface interface {
	Create(req *models.DBTenantLogin) error
	GetDetailsByEmail(email string) (*models.DBTenantLogin, error)
}

type TenantLoginRepository struct {
	DB *gorm.DB
}

func NewTenantLoginRepository(db *gorm.DB) (TenantLoginRepositoryInterface, error) {
	return &TenantLoginRepository{DB: db}, nil
}

func (tl *TenantLoginRepository) Create(req *models.DBTenantLogin) error {
	transaction := tl.DB.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	err := transaction.Create(&req)
	if err.Error != nil {
		return err.Error
	}
	transaction.Commit()
	return nil
}

func (tl *TenantLoginRepository) GetDetailsByEmail(email string) (*models.DBTenantLogin, error) {
	transaction := tl.DB.Begin()
	if transaction.Error != nil {
		return nil, transaction.Error
	}
	defer transaction.Rollback()
	var resp *models.DBTenantLogin
	err := transaction.Where("email = ? ", email).First(&resp)
	if err.Error != nil {
		return nil, err.Error
	}
	return resp, nil
}