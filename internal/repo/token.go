package repo

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

type TokenRepositoryInterface interface {
	CreateToken(token *models.DBToken) error
	UpdateToken(tenantid uuid.UUID) (string, error)
	ListTokens(tenantid uuid.UUID) (resp []*models.DBToken, err error)
	GetTenantUsingToken(token string) (*uuid.UUID, error)
	RevokeToken(token string) error
	VerifyToken(token string) (bool, string, error)
}

type TokenRepository struct {
	DB *gorm.DB
}

func NewTokenRepository(db *gorm.DB) (TokenRepositoryInterface, error) {
	return &TokenRepository{DB: db}, nil
}

func (to *TokenRepository) CreateToken(token *models.DBToken) error {
	transaction := to.DB.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	newToken := transaction.Create(&token)
	if newToken.Error != nil {
		return newToken.Error
	}
	transaction.Commit()
	return nil
}

func (to *TokenRepository) UpdateToken(tenantid uuid.UUID) (string, error) {
	transaction := to.DB.Begin()
	if transaction.Error != nil {
		return "", transaction.Error
	}
	defer transaction.Rollback()
	var tokenDetails *models.DBToken
	token := transaction.Model(&models.DBToken{}).Where("tenant_id = ?", tenantid).First(&tokenDetails)
	if token.Error != nil {
		return "", token.Error
	}
	newToken := uuid.New().String()
	if err := transaction.Model(&models.DBToken{}).Where("tenant_id = ?", tenantid).Updates(map[string]interface{}{
		"token": newToken,
	}).Error; err != nil {
		return "", err
	}
	transaction.Commit()
	return newToken, nil

}

func (to *TokenRepository) ListTokens(tenantid uuid.UUID) (resp []*models.DBToken, err error) {
	transaction := to.DB.Begin()
	if transaction.Error != nil {
		return nil, transaction.Error
	}
	defer transaction.Rollback()
	token := transaction.Model(&models.DBToken{}).Where("tenant_id = ?", tenantid).Find(&resp)
	if token.Error != nil {
		return nil, token.Error
	}
	return resp, err
}

func (to *TokenRepository) GetTenantUsingToken(token string) (*uuid.UUID, error) {
	transaction := to.DB.Begin()
	if transaction.Error != nil {
		return nil, transaction.Error
	}
	defer transaction.Rollback()
	var tokenDetails *models.DBToken
	id := transaction.Model(&models.DBToken{}).Where("token = ?", token).First(&tokenDetails)
	if id.Error != nil {
		return nil, id.Error
	}
	return &tokenDetails.TenantId, nil
}

func (to *TokenRepository) RevokeToken(token string) error {
	transaction := to.DB.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	err := transaction.Model(&models.DBToken{}).Where("token = ?", token).Updates(map[string]interface{}{
		"is_active": false,
	})
	if err != nil {
		return err.Error
	}
	return nil
}

func (to *TokenRepository) VerifyToken(token string) (bool, string, error) {
	log.Println("inside the verify token")
	transaction := to.DB.Begin()
	if transaction.Error != nil {
		return false, "", transaction.Error
	}
	defer transaction.Rollback()
	var tokenDetails models.DBToken
	tokenErr := transaction.Model(&models.DBToken{}).Where("token = ?", token).First(&tokenDetails)
	if tokenErr.Error != nil {
		log.Println("the token details error: ", tokenErr.Error)
		return false, "", transaction.Error
	}
	log.Println("the token details: ", tokenDetails)
	if tokenDetails.ExpiresAt.Before(time.Now()) {
		return false, "", &models.ServiceResponse{
			Code:    404,
			Message: "Token expired please login again",
		}
	}
	if !tokenDetails.IsActive {
		return false, "", &models.ServiceResponse{
			Code:    401,
			Message: "Token has been revoked. Please authenticate again.",
		}
	}
	rerr := to.RevokeToken(token)
	if rerr != nil {
		return false, "", rerr
	}
	return true, tokenDetails.TenantId.String(), nil
}
