package repo

import (
	"github.com/google/uuid"
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

type ResetTokenRepository struct {
	DB *gorm.DB
}

type ResetTokenRepositoryInterface interface {
	Create(req *models.DBResetToken) (uuid.UUID, error)
	FindAllToken(userId uuid.UUID) (resp []*models.DBResetToken, err error)
	VerifyOTP(otp string) (bool, error)
}

func NewResetTokenRepository(db *gorm.DB) (ResetTokenRepositoryInterface, error) {
	return &ResetTokenRepository{DB: db}, nil
}

func (rt *ResetTokenRepository) Create(req *models.DBResetToken) (uuid.UUID, error) {
	transaction := rt.DB.Begin()
	if transaction.Error != nil {
		return uuid.Nil, transaction.Error
	}
	defer transaction.Rollback()
	create := transaction.Create(&req)
	if create.Error != nil {
		return uuid.Nil, create.Error
	}
	transaction.Commit()
	return req.Id, nil
}

func (rt *ResetTokenRepository) FindAllToken(userId uuid.UUID) (resp []*models.DBResetToken, err error) {
	transaction := rt.DB.Begin()
	if transaction.Error != nil {
		return nil, transaction.Error
	}
	defer transaction.Rollback()
	result := transaction.Model(models.DBResetToken{}).Where("user_id = ? ", userId).Find(&resp)
	if result.Error != nil {
		return nil, result.Error
	}
	return resp, nil
}

func (rt *ResetTokenRepository) VerifyOTP(otp string) (bool, error) {
	transaction := rt.DB.Begin()
	if transaction.Error != nil {
		return false, transaction.Error
	}
	defer transaction.Rollback()
	var otpCheck *models.DBResetToken
	result := transaction.Model(models.DBResetToken{}).Where("otp = ? ", otp).First(&otpCheck)
	if result.Error != nil {
		return false, result.Error
	}
	update := transaction.Model(models.DBResetToken{}).Where("otp = ? ", otp).Updates(map[string]interface{}{
		"is_active": false,
	})
	if update.Error != nil {
		return false, update.Error
	}
	return true, nil
}
