package repo

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

type TokenRepositoryInterface interface {
	CreateToken(token *models.DBToken) error
	UpdateLoginToken(tenantid uuid.UUID) (*uuid.UUID, error)
	ListTokens(tenantid uuid.UUID) (resp []*models.DBToken, err error)
	GetTenantUsingToken(token string) (*uuid.UUID, error)
	RevokeToken(token string) error
	VerifyToken(token string) (bool, string, error)
	GetTokenDetailsByName(name string) (*models.DBToken, error)
	GetTokenDetails(conditions *models.DBToken) (*models.DBToken, error)
	VerifyApplicationToken(token string) (bool, string, error)
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

func (to *TokenRepository) UpdateLoginToken(tenantid uuid.UUID) (*uuid.UUID, error) {
	transaction := to.DB.Begin()
	if transaction.Error != nil {
		return nil, transaction.Error
	}
	defer transaction.Rollback()
	var tokenDetails *models.DBToken
	token := transaction.Model(&models.DBToken{}).Where(&models.DBToken{TenantId: tenantid, ApplicationKey: false}).First(&tokenDetails)
	if token.Error != nil {
		return nil, token.Error
	}
	newToken := uuid.New()
	if err := transaction.Model(&models.DBToken{}).Where(&models.DBToken{Id: tokenDetails.Id}).Updates(map[string]interface{}{
		"id":         newToken,
		"created_at": time.Now(),
		"expires_at": time.Now().Add(24 * time.Hour),
	}).Error; err != nil {
		return nil, err
	}
	transaction.Commit()
	return &newToken, nil

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
	id := transaction.Model(&models.DBToken{}).Where("id = ?", uuid.MustParse(token)).First(&tokenDetails)
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
	err := transaction.Model(&models.DBToken{}).Where("id = ?", uuid.MustParse(token)).Updates(map[string]any{
		"is_active":  false,
		"revoked_at": time.Now(),
	})
	if err.Error != nil {
		return err.Error
	}
	transaction.Commit()
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
	tokenErr := transaction.Model(&models.DBToken{}).Where("id = ?", uuid.MustParse(token)).First(&tokenDetails)
	if tokenErr.Error != nil {
		if tokenErr.Error.Error() == "record not found" {
			return false, "", errors.New("record not found")
		} else {
			return false, "", transaction.Error
		}
	}
	log.Println("the token details: ", tokenDetails)
	if tokenDetails.ExpiresAt.Before(time.Now()) {
		rerr := to.RevokeToken(token)
		if rerr != nil {
			return false, "", rerr
		}
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
	if tokenDetails.ApplicationKey {
		return false, "", &models.ServiceResponse{
			Code:    423,
			Message: "cant use application key as tenant login token",
		}
	}
	return true, tokenDetails.TenantId.String(), nil
}

func (to *TokenRepository) GetTokenDetailsByName(name string) (*models.DBToken, error) {
	transaction := to.DB.Begin()
	if transaction.Error != nil {
		return nil, transaction.Error
	}
	defer transaction.Rollback()
	var tokenDetails *models.DBToken
	err := transaction.Model(&models.DBToken{}).Where("name = ?", name).First(&tokenDetails)
	if err.Error != nil {
		return nil, err.Error
	}
	return tokenDetails, nil
}

func (to *TokenRepository) GetTokenDetails(conditions *models.DBToken) (*models.DBToken, error) {
	transaction := to.DB.Begin()
	if transaction.Error != nil {
		return nil, transaction.Error
	}
	defer transaction.Rollback()
	var tokenDetails *models.DBToken
	err := transaction.Model(&models.DBToken{}).First(&tokenDetails, &conditions)
	if err.Error != nil {
		return nil, err.Error
	}
	return tokenDetails, nil
}

func (to *TokenRepository) VerifyApplicationToken(token string) (bool, string, error) {
	log.Println("inside the verify token")
	transaction := to.DB.Begin()
	if transaction.Error != nil {
		return false, "", transaction.Error
	}
	defer transaction.Rollback()
	var tokenDetails models.DBToken
	tokenErr := transaction.Model(&models.DBToken{}).Where("id = ?", uuid.MustParse(token)).First(&tokenDetails)
	if tokenErr.Error != nil {
		if tokenErr.Error.Error() == "record not found" {
			return false, "", errors.New("record not found")
		} else {
			return false, "", transaction.Error
		}
	}
	log.Println("the token details: ", tokenDetails)
	if tokenDetails.ExpiresAt.Before(time.Now()) {
		rerr := to.RevokeToken(token)
		if rerr != nil {
			return false, "", rerr
		}
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
	if !tokenDetails.ApplicationKey {
		return false, "", &models.ServiceResponse{
			Code:    423,
			Message: "cant use login token as the application key",
		}
	}
	fmt.Println("The token was verrified successfully")
	fmt.Println("teh tenantId : ", tokenDetails.TenantId)
	return true, tokenDetails.TenantId.String(), nil
}
