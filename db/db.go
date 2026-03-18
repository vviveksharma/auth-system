package db

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Get DB host from environment, default to localhost for local development
	// Strip newlines from env vars to prevent log injection (G706)
	stripNL := func(s string) string {
		return strings.Map(func(r rune) rune {
			if r == '\n' || r == '\r' {
				return -1
			}
			return r
		}, s)
	}

	dbHost := stripNL(os.Getenv("DB_HOST"))
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := stripNL(os.Getenv("DB_PORT"))
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
}
