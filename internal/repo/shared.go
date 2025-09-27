package repo

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
	reqmodels "github.com/vviveksharma/auth/internal/models"
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

type SharedRepo struct {
	RoleRepo      RoleRepositoryInterface
	RouteRoleRepo RouteRoleRepositoryInterface
	DB            *gorm.DB
}

type SharedRepoInterface interface {
	CreateCustomRole(req *reqmodels.CreateCustomRole, tenantId uuid.UUID) error
	UpdateCustomRole(roleId uuid.UUID, tenantId uuid.UUID, addPermissions []reqmodels.Permission, removePermissions []reqmodels.Permission) error
	DeleteCustomRole(roleId uuid.UUID, tenantId uuid.UUID) error
}

func NewSharedRepository(db *gorm.DB) (SharedRepoInterface, error) {
	roleRepo, err := NewRoleRepository(db)
	if err != nil {
		return nil, errors.New("error from the shared reposistry with the role: " + err.Error())
	}
	routeRepo, err := NewRouteRoleRepository(db)
	if err != nil {
		return nil, errors.New("error from the shared reposistry with the route: " + err.Error())
	}
	return &SharedRepo{
		DB:            db,
		RoleRepo:      roleRepo,
		RouteRoleRepo: routeRepo}, nil
}

func (s *SharedRepo) CreateCustomRole(req *reqmodels.CreateCustomRole, tenantId uuid.UUID) error {
	transaction := s.DB.Begin()
	if transaction.Error != nil {
		fmt.Printf("Failed to begin transaction in FindRoleId: %v\n", transaction.Error)
		return transaction.Error
	}
	defer transaction.Rollback()
	roleId := uuid.New()
	// Create a role
	err := s.RoleRepo.CreateRole(&models.DBRoles{
		Role:     req.Name,
		TenantId: tenantId,
		RoleId:   roleId,
		RoleType: "custom",
		Status:   true,
	})
	if err != nil {
		return fmt.Errorf("error while creating a entry in the role database while creating the custom role: %s", err.Error())
	}
	// Generate the permissions
	fmt.Println("the requestBody: ", req)
	// Validate that req.Permissions is not nil/empty
	if req.Permissions == nil {
		return fmt.Errorf("permissions cannot be nil")
	}

	permissions := reqmodels.RoleData{
		RoleInfo: reqmodels.RoleInfo{
			Name:        req.Name,
			DisplayName: req.DisplayName,
			Description: req.Description,
			RoleType:    "custom",
			Priority:    50,
			IsSystem:    false,
		},
		Permissions: req.Permissions,
	}

	// Debug: Print the permissions to see what's being stored
	fmt.Printf("Creating role with permissions: %+v\n", permissions.Permissions)

	permissionsByte, err := json.Marshal(permissions)
	if err != nil {
		return fmt.Errorf("error while converting the permission to the string: %s", err.Error())
	}

	// Create a single RoleRouteEntry
	err = s.RouteRoleRepo.Create(&models.DBRouteRole{
		RoleName:    req.Name,
		TenantId:    tenantId,
		RoleId:      roleId,
		Permissions: string(permissionsByte),
		Routes:      pq.StringArray{"/temp"},
	})
	if err != nil {
		return fmt.Errorf("error while creating a entry in the role-route database while creating the custom role: %s", err.Error())
	}
	transaction.Commit()
	return nil
}

func (s *SharedRepo) UpdateCustomRole(roleId uuid.UUID, tenantId uuid.UUID, addPermissions []reqmodels.Permission, removePermissions []reqmodels.Permission) error {
	transaction := s.DB.Begin()
	if transaction.Error != nil {
		fmt.Printf("Failed to begin transaction in FindRoleId: %v\n", transaction.Error)
		return transaction.Error
	}
	permissions, err := s.RouteRoleRepo.GetRoleRouteMapping(roleId.String())
	if err != nil {
		return err
	}
	var roleData reqmodels.RoleData
	err = json.Unmarshal([]byte(permissions.Permissions), &roleData)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	fmt.Println("Before updating the permissions: ", roleData.Permissions)

	if len(removePermissions) > 0 {
		removedCount := s.removePermissionsWithLogging(roleData.Permissions, removePermissions)
		roleData.Permissions = removePermissionsFromRole(roleData.Permissions, removePermissions)
		fmt.Printf("Removed %d permissions from role %s\n", removedCount, roleId.String())
	}

	if len(addPermissions) > 0 {
		addedCount := len(addPermissions) - s.countExistingPermissions(roleData.Permissions, addPermissions)
		roleData.Permissions = addPermissionsToRole(roleData.Permissions, addPermissions)
		fmt.Printf("Added %d new permissions to role %s\n", addedCount, roleId.String())
	}
	defer transaction.Rollback()

	updatedPermissionsByte, err := json.Marshal(roleData)
	if err != nil {
		transaction.Rollback()
		return fmt.Errorf("error marshaling updated permissions: %v", err)
	}

	fmt.Println("After updating the permissions: ", roleData.Permissions)
	update := transaction.Model(&models.DBRouteRole{}).Where("role_id = ? AND tenant_id = ?", roleId, tenantId).Updates(map[string]interface{}{
		"permissions": string(updatedPermissionsByte),
	})
	if update.Error != nil {
		return fmt.Errorf("error while updating the route-rolemapping: %s", update.Error)
	}
	transaction.Commit()
	return nil
}

func (s *SharedRepo) DeleteCustomRole(roleId uuid.UUID, tenantId uuid.UUID) error {
	transaction := s.DB.Begin()
	if transaction.Error != nil {
		fmt.Printf("Failed to begin transaction in FindRoleId: %v\n", transaction.Error)
		return transaction.Error
	}
	defer transaction.Rollback()
	role := transaction.Model(models.DBRoles{}).Where("role_id = ? AND tenant_id = ? ", roleId, tenantId).Delete(models.DBRoles{
		RoleId:   roleId,
		TenantId: tenantId,
	})
	if role.Error != nil {
		return role.Error
	}
	route := transaction.Model(models.DBRouteRole{}).Where("role_id = ? AND tenant_id = ?", roleId, tenantId).Delete(models.DBRouteRole{
		RoleId:   roleId,
		TenantId: tenantId,
	})
	if route.Error != nil {
		return route.Error
	}
	transaction.Commit()
	return nil
}

func (s *SharedRepo) countExistingPermissions(existing []reqmodels.Permission, toAdd []reqmodels.Permission) int {
	count := 0
	for _, newPerm := range toAdd {
		if permissionExists(existing, newPerm) {
			count++
		}
	}
	return count
}

// Helper for logging removed permissions
func (s *SharedRepo) removePermissionsWithLogging(existing []reqmodels.Permission, toRemove []reqmodels.Permission) int {
	count := 0
	for _, existing := range existing {
		for _, remove := range toRemove {
			if permissionsMatch(existing, remove) {
				fmt.Printf("Removing permission: %s %v\n", remove.Route, remove.Methods)
				count++
				break
			}
		}
	}
	return count
}

func permissionsMatch(perm1, perm2 reqmodels.Permission) bool {
	if perm1.Route != perm2.Route {
		return false
	}

	// Compare methods (order independent)
	if len(perm1.Methods) != len(perm2.Methods) {
		return false
	}

	methodMap := make(map[string]bool)
	for _, method := range perm1.Methods {
		methodMap[method] = true
	}

	for _, method := range perm2.Methods {
		if !methodMap[method] {
			return false
		}
	}

	return true
}

func permissionExists(permissions []reqmodels.Permission, target reqmodels.Permission) bool {
	for _, perm := range permissions {
		if permissionsMatch(perm, target) {
			return true
		}
	}
	return false
}

func removePermissionsFromRole(existing []reqmodels.Permission, toRemove []reqmodels.Permission) []reqmodels.Permission {
	var result []reqmodels.Permission
	fmt.Println("Before updating: ", existing)
	for _, existingPerm := range existing {
		shouldKeep := true

		for _, removePerm := range toRemove {
			if permissionsMatch(existingPerm, removePerm) {
				shouldKeep = false
				break
			}
		}

		if shouldKeep {
			result = append(result, existingPerm)
		}
	}
	fmt.Println("After updating: ", result)
	return result
}

func addPermissionsToRole(existing []reqmodels.Permission, toAdd []reqmodels.Permission) []reqmodels.Permission {
	result := make([]reqmodels.Permission, len(existing))
	copy(result, existing)
	fmt.Println("the existing add role from before: ", existing)
	for _, newPerm := range toAdd {
		if !permissionExists(result, newPerm) {
			result = append(result, newPerm)
		}
	}
	fmt.Println("the result add role from after: ", result)
	return result
}
