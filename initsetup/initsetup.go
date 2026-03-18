package initsetup

import (
	_ "embed"
	"log"

	"github.com/vviveksharma/auth/db"
)

//go:embed seed_system_roles.sql
var seedSQL string

// systemRoleCount is the total number of gr.* roles defined in seed_system_roles.sql.
// Update this constant whenever a role is added or removed from that file.
const systemRoleCount = 10

// InitSetup ensures all gr.* system roles and their route-role mappings are present.
// It uses a single COUNT check to decide whether seeding is needed, then executes
// the embedded SQL in one transaction. All INSERTs use ON CONFLICT DO NOTHING so
// this is safe to run on every startup.
func InitSetup() {
	sqlDB, err := db.DB.DB()
	if err != nil {
		log.Fatal("initsetup: failed to get underlying sql.DB: ", err)
	}

	var count int
	if err := sqlDB.QueryRow(
		`SELECT COUNT(*) FROM role_tbl WHERE role LIKE 'gr.%'`,
	).Scan(&count); err != nil {
		log.Fatal("initsetup: failed to count system roles: ", err)
	}

	if count >= systemRoleCount {
		log.Printf("✅ initsetup: %d system roles already present — skipping seed", count)
		return
	}

	log.Printf("🌱 initsetup: found %d/%d system roles — running seed", count, systemRoleCount)

	tx, err := sqlDB.Begin()
	if err != nil {
		log.Fatal("initsetup: failed to begin transaction: ", err)
	}

	if _, err := tx.Exec(seedSQL); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Printf("⚠️  initsetup: rollback failed: %v", rbErr)
		}
		log.Fatal("initsetup: seed SQL failed: ", err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatal("initsetup: failed to commit seed transaction: ", err)
	}

	log.Println("✅ initsetup: system roles seeded successfully")
}
