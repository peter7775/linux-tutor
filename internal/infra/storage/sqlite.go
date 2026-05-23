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
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS progress (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		correct INTEGER NOT NULL DEFAULT 0,
		wrong INTEGER NOT NULL DEFAULT 0
	);`)
	if err != nil { return err }
	var n int
	if err := db.QueryRow(`SELECT COUNT(*) FROM progress`).Scan(&n); err != nil { return err }
	if n == 0 {
		_, err = db.Exec(`INSERT INTO progress(correct, wrong) VALUES (0,0)`)
	}
	return err
}
