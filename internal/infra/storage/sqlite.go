package storage

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
)

func Open(path string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return nil, err
	}
	return sql.Open("sqlite", path)
}
func Migrate(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS progress (id INTEGER PRIMARY KEY CHECK (id=1), correct INTEGER NOT NULL DEFAULT 0, wrong INTEGER NOT NULL DEFAULT 0); CREATE TABLE IF NOT EXISTS attempts (id INTEGER PRIMARY KEY AUTOINCREMENT, topic_code TEXT, prompt TEXT, answer TEXT, notes TEXT, score_delta INTEGER, created_at TEXT); INSERT OR IGNORE INTO progress(id, correct, wrong) VALUES (1,0,0);`)
	return err
}
