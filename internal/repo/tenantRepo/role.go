package tenantrepo

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	dbmodels "github.com/vviveksharma/auth/internal/models"
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

type TenantRoleRepositoryInterface interface {
	ListRoles(tenantId uuid.UUID, page int, pageSize int, status string, roleType string) ([]*models.DBRoles, int64, error)
	GetPermissions(tenantId uuid.UUID, roleId uuid.UUID) (*dbmodels.RoleData, error)
	UpdateRolePermissions(tenantId uuid.UUID, roleId uuid.UUID, permissons []dbmodels.Permission) error
}

type TenantRoleRepository struct {
	DB *gorm.DB
}

func NewTenantRoleRepository(db *gorm.DB) (TenantRoleRepositoryInterface, error) {
	return &TenantRoleRepository{DB: db}, nil
}

func (tr *TenantRoleRepository) ListRoles(tenantId uuid.UUID, page int, pageSize int, status string, roleType string) ([]*models.DBRoles, int64, error) {
	var totalCount int64
	var roles []*models.DBRoles

	transaction := tr.DB.Begin()
	if transaction.Error != nil {
		return nil, 0, transaction.Error
	}
	defer transaction.Rollback()

	systemTenantId := uuid.MustParse(models.GetSystemTenantId())

	// Build the count query based on roleType
	var countQuery *gorm.DB

	if roleType == "all" {
		// For "all", count both custom roles (user's tenant) and default roles (system tenant)
		countQuery = tr.DB.Model(&models.DBRoles{}).Where(
			"(tenant_id = ? AND role_type = 'custom') OR (tenant_id = ? AND role_type = 'default')",
			tenantId, systemTenantId,
		)
	} else if roleType == "default" {
		// For default roles, use system tenant ID
		countQuery = tr.DB.Model(&models.DBRoles{}).Where("tenant_id = ?", systemTenantId).Where("role_type = ?", "default")
	} else {
		// For custom roles, use the provided tenant ID
		countQuery = tr.DB.Model(&models.DBRoles{}).Where("tenant_id = ?", tenantId).Where("role_type = ?", roleType)
	}

	// Add status filter only if not "all"
	if status != "all" {
		var is_active bool
		if status == "active" {
			is_active = true
		} else {
			is_active = false
		}
		countQuery = countQuery.Where("status = ?", is_active)
	}

	if err := countQuery.Count(&totalCount).Error; err != nil {
		log.Printf("Error counting the tenant-roles: %v", err)
		return nil, 0, fmt.Errorf("error counting the roles present for this tenant: %v", err)
	}
	log.Printf("Total roles found: %d", totalCount)

	offset := (page - 1) * pageSize

	// Build the fetch query based on roleType
	var fetchQuery *gorm.DB

	if roleType == "all" {
		// For "all", fetch both custom roles (user's tenant) and default roles (system tenant)
		fetchQuery = tr.DB.Model(&models.DBRoles{}).Where(
			"(tenant_id = ? AND role_type = 'custom') OR (tenant_id = ? AND role_type = 'default')",
			tenantId, systemTenantId,
		)
	} else if roleType == "default" {
		// For default roles, use system tenant ID
		fetchQuery = tr.DB.Model(&models.DBRoles{}).Where("tenant_id = ?", systemTenantId).Where("role_type = ?", "default")
	} else {
		// For custom roles, use the provided tenant ID
		fetchQuery = tr.DB.Model(&models.DBRoles{}).Where("tenant_id = ?", tenantId).Where("role_type = ?", roleType)
	}

	// Add status filter only if not "all"
	if status != "all" {
		var is_active bool
		if status == "active" {
			is_active = true
		} else {
			is_active = false
		}
		fetchQuery = fetchQuery.Where("status = ?", is_active)
	}

	if roleType == "all" {
		fetchQuery = fetchQuery.Order("CASE WHEN role_type = 'default' THEN 0 ELSE 1 END, created_at DESC")
	} else {
		fetchQuery = fetchQuery.Order("created_at DESC")
	}

	if err := fetchQuery.Limit(pageSize).Offset(offset).Find(&roles).Error; err != nil {
		log.Printf("Error fetching paginated roles for the tenant: %v", err)
		return nil, 0, fmt.Errorf("error fetching roles for this tenant: %w", err)
	}

	log.Printf("Successfully fetched %d roles for this tenant (page %d, pageSize %d, total %d, status: %s, roleType: %s)",
		len(roles), page, pageSize, totalCount, status, roleType)

	return roles, totalCount, nil
}

func (tr *TenantRoleRepository) GetPermissions(tenantId uuid.UUID, roleId uuid.UUID) (*dbmodels.RoleData, error) {
	log.Printf("GetPermissions called for tenantId: %s, roleId: %s", tenantId, roleId)

	var roleDetails models.DBRoles
	err := tr.DB.Model(&models.DBRoles{}).
		Where("tenant_id = ?", tenantId).
		Where("role_id = ?", roleId).
		First(&roleDetails).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("role not found for tenant %s with role_id %s", tenantId, roleId)
		}
		return nil, fmt.Errorf("error fetching role details: %w", err)
	}

	log.Printf("Found role: %s (type: %s)", roleDetails.Role, roleDetails.RoleType)

	var routeRoleDetails models.DBRouteRole
	err = tr.DB.Model(&models.DBRouteRole{}).
		Where("tenant_id = ?", tenantId).
		Where("role_id = ?", roleId).
		First(&routeRoleDetails).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("route-role mapping not found for role_id %s", roleId)
		}
		return nil, fmt.Errorf("error fetching route-role details: %w", err)
	}

	log.Printf("Found route-role mapping with %d routes", len(routeRoleDetails.Routes))

	roleData, err := dbmodels.ConvertDBData(routeRoleDetails.Permissions)
	if err != nil {
		return nil, fmt.Errorf("error converting role permissions: %w", err)
	}

	log.Printf("Successfully converted permissions for role: %s", roleDetails.Role)
	return roleData, nil
}

func (tr *TenantRoleRepository) UpdateRolePermissions(tenantId uuid.UUID, roleId uuid.UUID, permissons []dbmodels.Permission) error {
	transaction := tr.DB.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	var existingRouteRole models.DBRouteRole
	err := transaction.Where("role_id = ? AND tenant_id = ?", roleId, tenantId).
		First(&existingRouteRole).Error
	if err != nil {
		return err
	}
	newPermissionsJSON, err := json.Marshal(permissons)
	if err != nil {
		return err
	}
	updateResult := transaction.Model(&models.DBRouteRole{}).
		Where("role_id = ? AND tenant_id = ?", roleId, tenantId).
		Updates(map[string]interface{}{
			"permissions": newPermissionsJSON,
		})

	if updateResult.Error != nil {
		log.Printf("Error updating permissions for role %s: %v", roleId, updateResult.Error)
		return fmt.Errorf("failed to update permissions: %w", updateResult.Error)
	}
	return nil
}
