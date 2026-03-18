package repo

import (
	"time"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/internal/utils"
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

type ResetCredsRepositry struct {
	DB *gorm.DB
}

type ResetCredsRepositryInterface interface {
	Create(userId uuid.UUID, tenantId uuid.UUID) ([]string, error)
	InvalidateAll(userId uuid.UUID, tenantId uuid.UUID) error
	UpdateUsage(userId uuid.UUID, tenantId uuid.UUID, tokenId uuid.UUID) error
	ListAllCreds(userId uuid.UUID, tenantId uuid.UUID) ([]*models.DBResetCreds, error)
	FindByCreds(userId uuid.UUID, tenantId uuid.UUID, tokenId uuid.UUID) (*models.DBResetCreds, error)
}

func NewResetCredRepositry(db *gorm.DB) (ResetCredsRepositryInterface, error) {
	return &ResetCredsRepositry{DB: db}, nil
}

func (rc *ResetCredsRepositry) Create(userId uuid.UUID, tenantId uuid.UUID) ([]string, error) {
	transaction := rc.DB.Begin()
	if transaction.Error != nil {
		return nil, transaction.Error
	}
	var codes []string
	for range 10 {
		var newCreds *models.DBResetCreds
		newcode := uuid.New().String()
		codes = append(codes, newcode)
		hashedCode, salt, err := utils.GeneratePasswordHash(newcode, utils.DefaultParams)
		if err != nil {
			return nil, err
		}
		newCreds = &models.DBResetCreds{
			TenantId:  tenantId,
			UserId:    userId,
			Active:    true,
			CreatedAt: time.Now(),
			CodeHash:  hashedCode,
			Salt:      salt,
		}
		if err := transaction.Create(newCreds).Error; err != nil {
			transaction.Rollback()
			return nil, err
		}
	}
	if err := transaction.Commit().Error; err != nil {
		return nil, err
	}
	return codes, nil
}

func (rc *ResetCredsRepositry) InvalidateAll(userId uuid.UUID, tenantId uuid.UUID) error {
	transaction := rc.DB.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	update := transaction.Model(&models.DBResetCreds{}).Where("tenant_id = ? AND user_id =?", tenantId, userId).Updates(map[string]interface{}{
		"active": false,
	})
	if update.Error != nil {
		transaction.Rollback()
		return update.Error
	}
	if err := transaction.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (rc *ResetCredsRepositry) UpdateUsage(userId uuid.UUID, tenantId uuid.UUID, tokenId uuid.UUID) error {
	transaction := rc.DB.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	now := time.Now()
	update := transaction.Model(&models.DBResetCreds{}).Where("tenant_id = ? AND user_id =? AND id = ?", tenantId, userId, tokenId).Updates(map[string]interface{}{
		"active":  false,
		"used_at": &now,
	})
	if update.Error != nil {
		return update.Error
	}
	if err := transaction.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (rc *ResetCredsRepositry) ListAllCreds(userId uuid.UUID, tenantId uuid.UUID) ([]*models.DBResetCreds, error) {
	transaction := rc.DB.Begin()
	if transaction.Error != nil {
		return nil, transaction.Error
	}
	defer transaction.Rollback()
	var response []*models.DBResetCreds
	list := transaction.Where("tenant_id = ? AND user_id = ? AND active = ?", tenantId, userId, true).Find(&response)
	if list.Error != nil {
		return nil, list.Error
	}
	return response, nil
}

func (rc *ResetCredsRepositry) FindByCreds(userId uuid.UUID, tenantId uuid.UUID, tokenId uuid.UUID) (*models.DBResetCreds, error) {
	transaction := rc.DB.Begin()
	if transaction.Error != nil {
		return nil, transaction.Error
	}
	defer transaction.Rollback()
	var response *models.DBResetCreds
	list := transaction.Where("tenant_id = ? AND user_id = ? AND id = ?", tenantId, userId, tokenId).Find(&response)
	if list.Error != nil {
		return nil, list.Error
	}
	transaction.Commit()
	return response, nil
}
