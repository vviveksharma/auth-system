package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	defaultMigrationsDir = "./db/migrations/sql"
	migrationsTable      = "schema_migrations"
)

type migration struct {
	name string
	sql  string
}

func main() {
	_ = godotenv.Load(".env")

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	fs := flag.NewFlagSet("migrate", flag.ExitOnError)
	dir := fs.String("dir", defaultMigrationsDir, "path to the SQL migrations directory")
	steps := fs.Int("steps", 1, "number of migrations to roll back (0 = all) – only used by 'down'")
	if err := fs.Parse(os.Args[2:]); err != nil {
		log.Fatal(err)
	}

	db := connectDB()
	defer db.Close()

	ensureMigrationsTable(db)

	switch command {
	case "up":
		runUp(db, *dir)
	case "down":
		runDown(db, *dir, *steps)
	case "status":
		showStatus(db, *dir)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %q\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func connectDB() *sql.DB {
	host := envOrDefault("DB_HOST", "localhost")
	port := envOrDefault("DB_PORT", "26257")

	dsn := fmt.Sprintf("postgresql://root@%s:%s/defaultdb?sslmode=disable", host, port)
	log.Printf("Connecting to database at %s:%s …", host, port)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("sql.Open: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("db.Ping: %v", err)
	}
	log.Println("Connected.")
	return db
}

func ensureMigrationsTable(db *sql.DB) {
	query := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
    name       TEXT        PRIMARY KEY,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
);`, migrationsTable)
	if _, err := db.Exec(query); err != nil {
		log.Fatalf("Failed to create migrations table: %v", err)
	}
}

func runUp(db *sql.DB, dir string) {
	migrations := loadMigrations(dir, "up")
	applied := appliedMigrations(db)

	ran := 0
	for _, m := range migrations {
		if applied[m.name] {
			continue
		}
		log.Printf("Applying  %-50s …", m.name)
		if err := execMigration(db, m.sql); err != nil {
			log.Fatalf("FAILED: %v", err)
		}
		markApplied(db, m.name)
		log.Printf("OK")
		ran++
	}

	if ran == 0 {
		fmt.Println("Nothing to migrate – all migrations are already applied.")
	} else {
		fmt.Printf("Applied %d migration(s).\n", ran)
	}
}

func runDown(db *sql.DB, dir string, steps int) {
	migrations := loadMigrations(dir, "down")
	// reverse: newest first
	for i, j := 0, len(migrations)-1; i < j; i, j = i+1, j-1 {
		migrations[i], migrations[j] = migrations[j], migrations[i]
	}

	applied := appliedMigrations(db)

	rolled := 0
	for _, m := range migrations {
		if steps > 0 && rolled >= steps {
			break
		}
		if !applied[m.name] {
			continue
		}
		log.Printf("Rolling back %-48s …", m.name)
		if err := execMigration(db, m.sql); err != nil {
			log.Fatalf("FAILED: %v", err)
		}
		markReverted(db, m.name)
		log.Printf("OK")
		rolled++
	}

	if rolled == 0 {
		fmt.Println("Nothing to roll back.")
	} else {
		fmt.Printf("Rolled back %d migration(s).\n", rolled)
	}
}

func showStatus(db *sql.DB, dir string) {
	upMigrations := loadMigrations(dir, "up")
	applied := appliedMigrations(db)

	fmt.Printf("\n%-6s  %-50s  %s\n", "STATUS", "MIGRATION", "APPLIED AT")
	fmt.Println(strings.Repeat("-", 80))

	for _, m := range upMigrations {
		if applied[m.name] {
			fmt.Printf("%-6s  %-50s  %s\n", "up", m.name, appliedAt(db, m.name).Format(time.RFC3339))
		} else {
			fmt.Printf("%-6s  %-50s  %s\n", "down", m.name, "-")
		}
	}
	fmt.Println()
}

func loadMigrations(dir, direction string) []migration {
	pattern := filepath.Join(dir, "*."+direction+".sql")
	files, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatalf("glob %q: %v", pattern, err)
	}
	if len(files) == 0 {
		log.Fatalf("No *.%s.sql files found in %q", direction, dir)
	}
	sort.Strings(files)

	var result []migration
	for _, f := range files {
		base := filepath.Base(f)
		name := strings.TrimSuffix(base, "."+direction+".sql")
		content, err := os.ReadFile(f) // #nosec G304 -- CLI tool reading its own SQL files, no user input
		if err != nil {
			log.Fatalf("ReadFile %q: %v", f, err)
		}
		result = append(result, migration{name: name, sql: string(content)})
	}
	return result
}

func appliedMigrations(db *sql.DB) map[string]bool {
	rows, err := db.Query(fmt.Sprintf("SELECT name FROM %s", migrationsTable))
	if err != nil {
		log.Fatalf("query applied migrations: %v", err)
	}
	defer rows.Close()

	m := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatalf("scan: %v", err)
		}
		m[name] = true
	}
	return m
}

func execMigration(db *sql.DB, sqlContent string) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	if _, err := tx.Exec(sqlContent); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func markApplied(db *sql.DB, name string) {
	_, err := db.Exec(
		fmt.Sprintf("INSERT INTO %s (name) VALUES ($1) ON CONFLICT DO NOTHING", migrationsTable),
		name,
	)
	if err != nil {
		log.Fatalf("markApplied %q: %v", name, err)
	}
}

func markReverted(db *sql.DB, name string) {
	_, err := db.Exec(
		fmt.Sprintf("DELETE FROM %s WHERE name = $1", migrationsTable),
		name,
	)
	if err != nil {
		log.Fatalf("markReverted %q: %v", name, err)
	}
}

func appliedAt(db *sql.DB, name string) time.Time {
	var t time.Time
	_ = db.QueryRow(
		fmt.Sprintf("SELECT applied_at FROM %s WHERE name = $1", migrationsTable),
		name,
	).Scan(&t)
	return t
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func printUsage() {
	fmt.Print(`Usage:
  go run ./cmd/migrate <command> [flags]

Commands:
  up       Apply all pending migrations
  down     Roll back migrations (default: 1 step)
  status   Show applied vs pending migrations

Flags:
  -dir    string   Path to SQL migrations directory (default: ./db/migrations/sql)
  -steps  int      Number of migrations to roll back with 'down' (default: 1; 0 = all)

Examples:
  go run ./cmd/migrate up
  go run ./cmd/migrate down
  go run ./cmd/migrate down -steps 3
  go run ./cmd/migrate status
`)
}
