package migrations

import (
	"fmt"

	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

func AutoMigrator(DB *gorm.DB) {
	DB.AutoMigrate(models.DBUser{}, models.DBRoles{}, models.DBLogin{})
	fmt.Println("Migrations done!!!")
}
