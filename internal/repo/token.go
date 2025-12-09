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
	ListTokensPaginated(tenantid uuid.UUID, page, pageSize int, status string) ([]*models.DBToken, int64, error)
	ListTokens(tenantid uuid.UUID) (resp []*models.DBToken, err error)
	GetTenantUsingToken(token string) (*uuid.UUID, error)
	RevokeToken(token string) error
	VerifyToken(token string) (bool, string, error)
	GetTokenDetails(conditions models.DBToken) (*models.DBToken, error)
	VerifyApplicationToken(token string) (bool, string, error)
	GetTokenDetailsStatus(status string, tenantId uuid.UUID, page, pageSize int) ([]*models.DBToken, int64, error)
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

	var tokenDetails models.DBToken

	err := transaction.Model(&models.DBToken{}).
		Where("tenant_id = ?", tenantid).
		Where("application_key = ?", false).
		First(&tokenDetails).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("No active login token found for tenant: %s", tenantid)
			return nil, fmt.Errorf("no active login token found for tenant: %s", tenantid)
		}
		log.Printf("Error fetching token: %v", err)
		return nil, fmt.Errorf("error fetching token: %w", err)
	}
	log.Printf("Found existing token: ID=%s, Name=%s, TenantID=%s, ApplicationKey=%v",
		tokenDetails.Id, tokenDetails.Name, tokenDetails.TenantId, tokenDetails.ApplicationKey)

	newToken := uuid.New()

	log.Printf("Attempting to update token from %s to %s", tokenDetails.Id, newToken)
	updateResult := transaction.Model(&models.DBToken{}).
		Where("id = ?", tokenDetails.Id).
		Where("tenant_id = ?", tenantid).
		Updates(map[string]interface{}{
			"id":         newToken,
			"created_at": time.Now(),
			"expires_at": time.Now().Add(24 * time.Hour),
			"is_active":  true,
		})

	if updateResult.Error != nil {
		log.Printf("Error updating token: %v", updateResult.Error)
		return nil, fmt.Errorf("error updating token: %w", updateResult.Error)
	}
	if updateResult.RowsAffected == 0 {
		log.Printf("Warning: No rows affected during token update for ID=%s", tokenDetails.Id)
		return nil, fmt.Errorf("token update failed: no rows affected")
	}

	log.Printf("Successfully updated token. Rows affected: %d", updateResult.RowsAffected)
	if err := transaction.Commit().Error; err != nil {
		log.Printf("Error committing transaction: %v", err)
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	log.Printf("Transaction committed successfully. New token: %s", newToken)
	return &newToken, nil
}

func (to *TokenRepository) ListTokensPaginated(tenantid uuid.UUID, page, pageSize int, status string) ([]*models.DBToken, int64, error) {
	log.Printf("ListTokens called: tenantId=%s, page=%d, pageSize=%d", tenantid, page, pageSize)

	var totalCount int64
	var tokens []*models.DBToken
	var is_active bool

	if status == "active" {
		is_active = true
	} else {
		is_active = false
	}

	if err := to.DB.Model(&models.DBToken{}).
		Where("tenant_id = ?", tenantid).Where("is_active = ?", is_active).Where("application_key = ?", true).
		Count(&totalCount).Error; err != nil {
		log.Printf("Error counting tokens: %v", err)
		return nil, 0, fmt.Errorf("error counting tokens: %w", err)
	}

	log.Printf("Total tokens found: %d", totalCount)

	offset := (page - 1) * pageSize

	if err := to.DB.Model(&models.DBToken{}).
		Where("tenant_id = ?", tenantid).Where("is_active = ?", is_active).Where("application_key = ?", true).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&tokens).Error; err != nil {
		log.Printf("Error fetching paginated tokens: %v", err)
		return nil, 0, fmt.Errorf("error fetching tokens: %w", err)
	}

	log.Printf("Successfully fetched %d tokens (page %d, pageSize %d, total %d)",
		len(tokens), page, pageSize, totalCount)

	return tokens, totalCount, nil
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
	log.Println("the verify token: ", token)
	var tokenDetails models.DBToken
	tokenErr := transaction.Model(&models.DBToken{}).Where("id = ?", uuid.MustParse(token)).First(&tokenDetails)
	if tokenErr.Error != nil {
		if errors.Is(tokenErr.Error, gorm.ErrRecordNotFound) {
			return false, "", errors.New("record not found")
		} else {
			return false, "", tokenErr.Error
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
			Code:    409,
			Message: "cant use application key as tenant login token",
		}
	}
	return true, tokenDetails.TenantId.String(), nil
}

func (to *TokenRepository) GetTokenDetails(conditions models.DBToken) (*models.DBToken, error) {
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
			Code:    409,
			Message: "cant use login token as the application key",
		}
	}
	fmt.Println("The token was verrified successfully")
	fmt.Println("teh tenantId : ", tokenDetails.TenantId)
	return true, tokenDetails.TenantId.String(), nil
}

func (to *TokenRepository) GetTokenDetailsStatus(status string, tenantId uuid.UUID, page, pageSize int) ([]*models.DBToken, int64, error) {
	log.Printf("GetTokenDetailsStatus called: tenantId=%s, status=%s, page=%d, pageSize=%d", tenantId, status, page, pageSize)

	var totalCount int64
	var tokens []*models.DBToken

	var isActive bool
	switch status {
	case "true", "active", "1":
		isActive = true
	case "false", "inactive", "0":
		isActive = false
	default:
		return nil, 0, fmt.Errorf("invalid status value: %s. Use 'true'/'false' or 'active'/'inactive'", status)
	}

	baseQuery := to.DB.Model(&models.DBToken{}).
		Where("tenant_id = ?", tenantId).
		Where("is_active = ?", isActive)

	if err := baseQuery.Count(&totalCount).Error; err != nil {
		log.Printf("Error counting tokens: %v", err)
		return nil, 0, fmt.Errorf("error counting tokens: %w", err)
	}

	log.Printf("Total tokens found: %d", totalCount)

	offset := (page - 1) * pageSize

	query := to.DB.Model(&models.DBToken{}).
		Where("tenant_id = ?", tenantId).
		Where("is_active = ?", isActive).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset)

	if err := query.Find(&tokens).Error; err != nil {
		log.Printf("Error fetching paginated tokens: %v", err)
		return nil, 0, fmt.Errorf("error fetching tokens: %w", err)
	}

	log.Printf("Successfully fetched %d tokens (page %d, pageSize %d, total %d)",
		len(tokens), page, pageSize, totalCount)

	return tokens, totalCount, nil
}
