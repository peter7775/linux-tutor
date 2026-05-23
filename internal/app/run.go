package app

import (
	"database/sql"
	"fmt"
	"linux-tutor/internal/gui"
	"linux-tutor/internal/infra/storage"
	"linux-tutor/internal/terminal"
	"log/slog"
	"os"
)

func openDB() *sql.DB {
	db, err := storage.Open("data/linux-tutor.db")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := storage.Migrate(db); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return db
}
func RunGUI() {
	db := openDB()
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("close db", "err", err)
		}
	}()
	gui.Start(db)
}
func RunTUI() {
	db := openDB()
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("close db", "err", err)
		}
	}()
	if err := terminal.Start(db); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
