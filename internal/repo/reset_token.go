package repo

import (
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

type ResetTokenRepository struct {
	DB *gorm.DB
}

type ResetTokenRepositoryInterface interface {
	Create(req *models.DBResetToken) error
}

func NewResetTokenRepository(db *gorm.DB) (ResetTokenRepositoryInterface, error) {
	return &ResetTokenRepository{DB: db}, nil
}

func (rt *ResetTokenRepository) Create(req *models.DBResetToken) error {
	transaction := rt.DB.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	create := transaction.Create(&req)
	if create.Error != nil {
		return create.Error
	}
	transaction.Commit()
	return nil
}
