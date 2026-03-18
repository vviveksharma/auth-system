package migrations

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"sort"
	"strings"
)

//go:embed sql
var sqlFS embed.FS

const migrationsTable = "schema_migrations"

// RunUp applies every pending *.up.sql file that has not yet been recorded in
// schema_migrations. SQL files are embedded at compile time so no filesystem
// path is needed — the binary is fully self-contained in Docker or anywhere else.
// It is idempotent and safe to call on every server startup.
func RunUp(db *sql.DB) error {
	if err := ensureMigrationsTable(db); err != nil {
		return fmt.Errorf("ensure migrations table: %w", err)
	}

	files := loadEmbeddedFiles("up")
	applied := appliedSet(db)

	ran := 0
	for _, f := range files {
		if applied[f.name] {
			continue
		}
		log.Printf("⬆  Applying migration: %s", f.name)
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("begin tx for %s: %w", f.name, err)
		}
		// Strip any BEGIN/COMMIT/ROLLBACK the file may contain — the runner
		// manages the transaction itself, and CockroachDB rejects nested BEGIN.
		if _, err := tx.Exec(stripTransactionWrappers(f.sql)); err != nil {
			tx.Rollback()
			return fmt.Errorf("execute migration %s: %w", f.name, err)
		}
		if _, err := tx.Exec(
			fmt.Sprintf(`INSERT INTO %s (name) VALUES ($1)`, migrationsTable),
			f.name,
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("mark migration %s applied: %w", f.name, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %s: %w", f.name, err)
		}
		log.Printf("✅ Migration applied: %s", f.name)
		ran++
	}

	if ran == 0 {
		log.Println("✅ All migrations already applied — nothing to do.")
	} else {
		log.Printf("✅ Applied %d migration(s) successfully.", ran)
	}
	return nil
}

type migFile struct {
	name string
	sql  string
}

func ensureMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			name       TEXT        PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)`, migrationsTable))
	return err
}

func loadEmbeddedFiles(direction string) []migFile {
	entries, err := sqlFS.ReadDir("sql")
	if err != nil {
		log.Fatalf("Cannot read embedded migrations: %v", err)
	}

	suffix := "." + direction + ".sql"
	var files []migFile
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), suffix) {
			continue
		}
		content, err := sqlFS.ReadFile("sql/" + e.Name())
		if err != nil {
			log.Fatalf("Cannot read embedded migration file %s: %v", e.Name(), err)
		}
		// Store without direction suffix so up/down share the same key in schema_migrations
		name := strings.TrimSuffix(e.Name(), suffix)
		files = append(files, migFile{name: name, sql: string(content)})
	}

	sort.Slice(files, func(i, j int) bool { return files[i].name < files[j].name })
	return files
}

// stripTransactionWrappers removes bare BEGIN, COMMIT, and ROLLBACK statements
// from SQL text so migration files that include them for use with external
// tools (psql, migrate CLI) can be run safely inside the runner's own tx.
func stripTransactionWrappers(sql string) string {
	lines := strings.Split(sql, "\n")
	kept := lines[:0]
	for _, line := range lines {
		trimmed := strings.TrimSpace(strings.ToUpper(line))
		if trimmed == "BEGIN" || trimmed == "BEGIN;" ||
			trimmed == "COMMIT" || trimmed == "COMMIT;" ||
			trimmed == "ROLLBACK" || trimmed == "ROLLBACK;" {
			continue
		}
		kept = append(kept, line)
	}
	return strings.Join(kept, "\n")
}

func appliedSet(db *sql.DB) map[string]bool {
	rows, err := db.Query(fmt.Sprintf(`SELECT name FROM %s`, migrationsTable))
	if err != nil {
		log.Fatalf("Query applied migrations: %v", err)
	}
	defer rows.Close()

	applied := map[string]bool{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatalf("Scan migration name: %v", err)
		}
		applied[name] = true
	}
	return applied
}
