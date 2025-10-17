package initsetup

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/db"
	"github.com/vviveksharma/auth/internal/repo"
	"github.com/vviveksharma/auth/internal/utils"
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

// Predefined role IDs
var (
	AdminId       = uuid.MustParse("f47ac10b-58cc-4372-a567-0e02b2c3d479")
	UserId        = uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	GuestId       = uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	ModeratorId   = uuid.MustParse("1b9d6bcd-bbfd-4b2d-9b5d-ab8dfbbd4bed")
	TenantId      = uuid.MustParse("dae760ab-0a7f-4cbd-8603-def85ad8e430")
	requiredRoles = []models.DBRoles{
		{Role: "admin", RoleId: AdminId, RoleType: "default", TenantId: TenantId, DisplayName: "Administrator"},
		{Role: "user", RoleId: UserId, RoleType: "default", TenantId: TenantId, DisplayName: "Content Moderator"},
		{Role: "guest", RoleId: GuestId, RoleType: "default", TenantId: TenantId, DisplayName: "Standard User"},
		{Role: "moderator", RoleId: ModeratorId, RoleType: "default", TenantId: TenantId, DisplayName: "Guest User"},
	}
)

var (
	requiredRolesRoutes = []models.DBRouteRole{
		{TenantId: uuid.MustParse("dae760ab-0a7f-4cbd-8603-def85ad8e430"), RoleId: uuid.MustParse("f47ac10b-58cc-4372-a567-0e02b2c3d479"), RoleName: "admin"},
		{TenantId: uuid.MustParse("dae760ab-0a7f-4cbd-8603-def85ad8e430"), RoleId: uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"), RoleName: "user"},
		{TenantId: uuid.MustParse("dae760ab-0a7f-4cbd-8603-def85ad8e430"), RoleId: uuid.MustParse("1b9d6bcd-bbfd-4b2d-9b5d-ab8dfbbd4bed"), RoleName: "moderator"},
		{TenantId: uuid.MustParse("dae760ab-0a7f-4cbd-8603-def85ad8e430"), RoleId: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"), RoleName: "guest"},
	}
)

func InitSetup() {
	db := db.DB
	exist, err := CheckRolesExist(db)
	if err != nil {
		log.Fatal("error checking roles existence: ", err)
	}

	var count int64
	err = db.Model(&models.DBRouteRole{}).Count(&count).Error
	if err != nil {
		log.Fatalln("error while getting the route and role: " + err.Error())
	}
	if count >= int64(len(requiredRolesRoutes)) {
		log.Println("Routes and role mapping already present")
		return
	} else {
		err = UpdateRoleRoutePermissions()
		if err != nil {
			log.Fatalln("error while updating the route and role mapping: " + err.Error())
		}
	}

	if exist {
		log.Println("roles already exist - skipping creation")
		return
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		for _, role := range requiredRoles {
			if err := tx.Create(&role).Error; err != nil {
				return fmt.Errorf("failed to create role %s: %w", role.Role, err)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal("error creating roles: ", err)
	}

	log.Println("roles created successfully and mapped to respective routes")
}

func CheckRolesExist(db *gorm.DB) (bool, error) {
	// Count the roles
	var count int64
	err := db.Model(&models.DBRoles{}).Count(&count).Error
	if err != nil {
		return false, err
	}

	leng := len(requiredRoles)

	if count >= int64(leng) {
		return true, nil
	}

	for _, role := range requiredRoles {
		var existing models.DBRoles
		err := db.Where("role_id = ?", role.RoleId).First(&existing).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return false, nil
			}
			return false, err
		}
	}
	return true, nil
}

func UpdateRoleRoutePermissions() error {
	log.Println("Inside the update permissions")
	roleRoute, err := repo.NewRouteRoleRepository(db.DB)
	if err != nil {
		log.Fatalf("error while connecting to the role route repositery: %s", err.Error())
		return err
	}
	for _, rr := range requiredRolesRoutes {
		permissions, err := utils.ReadPermissionFile(rr.RoleName)
		if err != nil {
			log.Fatalln("error while reding the permissions file: " + err.Error())
			return err
		}
		err = roleRoute.Create(&models.DBRouteRole{
			RoleName:    rr.RoleName,
			TenantId:    rr.TenantId,
			RoleId:      rr.RoleId,
			Permissions: permissions,
		})
		if err != nil {
			log.Fatalln("error while creating the roleRoute permission entry: " + err.Error())
			return err
		}
	}
	return nil
}
