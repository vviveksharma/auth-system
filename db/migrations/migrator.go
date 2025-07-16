package migrations

import (
	"fmt"
	"log"

	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

func AutoMigrator(DB *gorm.DB) {
	err := DB.AutoMigrate(models.DBUser{}, models.DBRoles{}, models.DBLogin{}, models.DBTenant{}, models.DBToken{}, models.DBTenantLogin{})
	if err != nil {
		log.Fatalln("error while migrating the tables: ", err.Error())
	}
	fmt.Println("Migrations done!!!")
}
