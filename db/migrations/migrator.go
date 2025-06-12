package migrations

import (
	"fmt"

	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

func AutoMigrator(DB *gorm.DB) {
	DB.AutoMigrate(models.DBUser{})
	fmt.Println("Migrations done!!!")
}
