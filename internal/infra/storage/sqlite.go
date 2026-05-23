package storage

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
)

func Open(path string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil { return nil, err }
	return sql.Open("sqlite", path)
}

func Migrate(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS progress (id INTEGER PRIMARY KEY CHECK (id=1), correct INTEGER NOT NULL DEFAULT 0, wrong INTEGER NOT NULL DEFAULT 0);`,
		`INSERT OR IGNORE INTO progress(id, correct, wrong) VALUES (1,0,0);`,
	}
	for _, s := range stmts { if _, err := db.Exec(s); err != nil { return err } }
	return nil
}
