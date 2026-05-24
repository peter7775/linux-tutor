package app

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"linux-tutor/internal/gui"
	"linux-tutor/internal/terminal"

	_ "github.com/glebarez/go-sqlite"
)

func RunGUI() error {
	db, err := openDB()
	if err != nil {
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("close db", "err", err)
		}
	}()
	gui.Start(db)
	return nil
}

func RunTUI() error {
	db, err := openDB()
	if err != nil {
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("close db", "err", err)
		}
	}()
	terminal.Start(db)
	return nil
}

func openDB() (*sql.DB, error) {
	path := filepath.Join("data", "linux-tutor.db")
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return nil, fmt.Errorf("create db dir: %w", err)
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	if err := db.Ping(); err != nil {
		if err := db.Close(); err != nil {
			slog.Error("close db", "err", err)
		}
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}
	return db, nil
}
