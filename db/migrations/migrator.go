package migrations

import (
	"fmt"
	"log"

	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

func AutoMigrator(DB *gorm.DB) {
	// Use GORM's Migrator to check if tables exist before migrating
	migrator := DB.Migrator()

	// List of models to migrate
	modelsToMigrate := []interface{}{
		&models.DBUser{},
		&models.DBRoles{},
		&models.DBLogin{},
		&models.DBTenant{},
		&models.DBToken{},
		&models.DBTenantLogin{},
		&models.DBRouteRole{},
		&models.DBResetToken{},
		&models.DBMessage{},
	}

	// Migrate each model individually with error handling
	for _, model := range modelsToMigrate {
		if !migrator.HasTable(model) {
			log.Printf("Creating table for %T...", model)
			if err := DB.AutoMigrate(model); err != nil {
				log.Printf("Warning: Error migrating %T: %v", model, err)
			}
		} else {
			// Table exists, just update columns if needed
			if err := DB.AutoMigrate(model); err != nil {
				// Ignore "already exists" errors
				if err.Error() != "ERROR: relation \"defaultdb.public.user_tbl\" already exists (SQLSTATE 42P07)" {
					log.Printf("Warning: Error updating table for %T: %v", model, err)
				}
			}
		}
	}

	fmt.Println("Migrations done!!!")
}
