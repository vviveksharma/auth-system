package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/vviveksharma/auth/db/migrations"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Get DB host from environment, default to localhost for local development
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "26257"
	}

	dsn := fmt.Sprintf("postgresql://root@%s:%s/defaultdb?sslmode=disable", dbHost, dbPort)
	log.Printf("Connecting to database at %s:%s...", dbHost, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to CockroachDB: ", err)
	}
	DB = db
	sqlDb, err := DB.DB()
	if err != nil {
		fmt.Printf("failed to get database connection: %v", err)
	}
	// Setting the scallable options
	sqlDb.SetMaxIdleConns(25)
	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetConnMaxLifetime(time.Hour)
	sqlDb.SetConnMaxIdleTime(10 * time.Minute)
	fmt.Println("The database ping returned: ", sqlDb.Ping())

	// Making migrations
	migrations.AutoMigrator(DB)
}
