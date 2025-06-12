package db

import (
	"fmt"
	"log"

	"github.com/vviveksharma/auth/db/migrations"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "postgresql://root@localhost:26257/defaultdb?sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to CockroachDB: ", err)
	}
	DB = db
	sqlDb, err := DB.DB()
	if err != nil {
        fmt.Printf("failed to get database connection: %v", err)
    }
    fmt.Println("The database ping returned: ",sqlDb.Ping())

	// Making migrations
	migrations.AutoMigrator(DB)
}
