package repo

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

type TenantRepositoryInterface interface {
	CreateTenant(tenant *models.DBTenant) error
	GetUserByEmail(email string) (tenantDetails *models.DBTenant, err error)
	UpdateTenatDetailsPassword(tenantId string, password string) error
	GetTenantDetails(conditions *models.DBTenant) (*models.DBTenant, error)
	DeleteTenant(tenantId uuid.UUID) error
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

func (t *TenantRepository) UpdateTenatDetailsPassword(tenantId string, password string) error {
	transaction := t.DB.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	updateErr := transaction.Model(&models.DBTenant{}).Where("id = ? ", uuid.MustParse(tenantId)).Updates(map[string]any{
		"password": password,
	})
	if updateErr.Error != nil {
		return updateErr.Error
	}
	transaction.Commit()
	return nil
}

func (t *TenantRepository) GetTenantDetails(conditions *models.DBTenant) (*models.DBTenant, error) {
	transaction := t.DB.Begin()
	if transaction.Error != nil {
		return nil, transaction.Error
	}
	defer transaction.Rollback()
	var tenantDetails models.DBTenant
	findErr := transaction.Model(&models.DBTenant{}).First(&tenantDetails, &conditions)
	if findErr.Error != nil {
		return nil, findErr.Error
	}
	return &tenantDetails, nil
}

func (t *TenantRepository) DeleteTenant(tenantId uuid.UUID) error {
	log.Printf("Starting cascading delete for tenant: %s", tenantId)

	transaction := t.DB.Begin()
	if transaction.Error != nil {
		log.Printf("Error starting transaction: %v", transaction.Error)
		return transaction.Error
	}
	defer transaction.Rollback()

	// Step 1: Delete all tokens associated with this tenant
	log.Printf("Deleting tokens for tenant: %s", tenantId)
	deleteTokens := transaction.Where("tenant_id = ?", tenantId).Delete(&models.DBToken{})
	if deleteTokens.Error != nil {
		log.Printf("Error deleting tokens: %v", deleteTokens.Error)
		return fmt.Errorf("failed to delete tokens: %w", deleteTokens.Error)
	}
	log.Printf("Deleted %d tokens", deleteTokens.RowsAffected)

	// Step 2: Delete all login records associated with this tenant
	log.Printf("Deleting login records for tenant: %s", tenantId)
	deleteLogins := transaction.Where("tenant_id = ?", tenantId).Delete(&models.DBLogin{})
	if deleteLogins.Error != nil {
		log.Printf("Error deleting login records: %v", deleteLogins.Error)
		return fmt.Errorf("failed to delete login records: %w", deleteLogins.Error)
	}
	log.Printf("Deleted %d login records", deleteLogins.RowsAffected)

	// Step 3: Delete all tenant login records
	log.Printf("Deleting tenant login records for tenant: %s", tenantId)
	deleteTenantLogins := transaction.Where("tenant_id = ?", tenantId).Delete(&models.DBTenantLogin{})
	if deleteTenantLogins.Error != nil {
		log.Printf("Error deleting tenant login records: %v", deleteTenantLogins.Error)
		return fmt.Errorf("failed to delete tenant login records: %w", deleteTenantLogins.Error)
	}
	log.Printf("Deleted %d tenant login records", deleteTenantLogins.RowsAffected)

	// Step 4: Delete all route-role mappings for roles in this tenant
	log.Printf("Deleting route-role mappings for tenant: %s", tenantId)
	deleteRouteRoles := transaction.Where("tenant_id = ?", tenantId).Delete(&models.DBRouteRole{})
	if deleteRouteRoles.Error != nil {
		log.Printf("Error deleting route-role mappings: %v", deleteRouteRoles.Error)
		return fmt.Errorf("failed to delete route-role mappings: %w", deleteRouteRoles.Error)
	}
	log.Printf("Deleted %d route-role mappings", deleteRouteRoles.RowsAffected)

	// Step 5: Delete all roles associated with this tenant
	log.Printf("Deleting roles for tenant: %s", tenantId)
	deleteRoles := transaction.Where("tenant_id = ?", tenantId).Delete(&models.DBRoles{})
	if deleteRoles.Error != nil {
		log.Printf("Error deleting roles: %v", deleteRoles.Error)
		return fmt.Errorf("failed to delete roles: %w", deleteRoles.Error)
	}
	log.Printf("Deleted %d roles", deleteRoles.RowsAffected)

	// Step 6: Delete all users associated with this tenant
	log.Printf("Deleting users for tenant: %s", tenantId)
	deleteUsers := transaction.Where("tenant_id = ?", tenantId).Delete(&models.DBUser{})
	if deleteUsers.Error != nil {
		log.Printf("Error deleting users: %v", deleteUsers.Error)
		return fmt.Errorf("failed to delete users: %w", deleteUsers.Error)
	}
	log.Printf("Deleted %d users", deleteUsers.RowsAffected)

	// Step 7: Delete password reset tokens associated with this tenant
	log.Printf("Deleting reset tokens for tenant: %s", tenantId)
	deleteResetTokens := transaction.Where("tenant_id = ?", tenantId).Delete(&models.DBResetToken{})
	if deleteResetTokens.Error != nil {
		log.Printf("Error deleting reset tokens: %v", deleteResetTokens.Error)
		return fmt.Errorf("failed to delete reset tokens: %w", deleteResetTokens.Error)
	}
	log.Printf("Deleted %d reset tokens", deleteResetTokens.RowsAffected)

	// Step 8: Finally, delete the tenant itself
	log.Printf("Deleting tenant: %s", tenantId)
	deleteTenant := transaction.Where("id = ?", tenantId).Delete(&models.DBTenant{})
	if deleteTenant.Error != nil {
		log.Printf("Error deleting tenant: %v", deleteTenant.Error)
		return fmt.Errorf("failed to delete tenant: %w", deleteTenant.Error)
	}

	if deleteTenant.RowsAffected == 0 {
		log.Printf("No tenant found with ID: %s", tenantId)
		return fmt.Errorf("tenant not found with ID: %s", tenantId)
	}
	log.Printf("Deleted tenant with ID: %s", tenantId)

	// Commit the transaction
	if err := transaction.Commit().Error; err != nil {
		log.Printf("Error committing transaction: %v", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Successfully completed cascading delete for tenant: %s", tenantId)
	log.Printf("Summary - Tokens: %d, Logins: %d, TenantLogins: %d, RouteRoles: %d, Roles: %d, Users: %d, ResetTokens: %d",
		deleteTokens.RowsAffected,
		deleteLogins.RowsAffected,
		deleteTenantLogins.RowsAffected,
		deleteRouteRoles.RowsAffected,
		deleteRoles.RowsAffected,
		deleteUsers.RowsAffected,
		deleteResetTokens.RowsAffected)

	return nil
}
