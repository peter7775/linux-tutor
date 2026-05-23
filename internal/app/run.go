package app

import (
	"fmt"
	"linux-tutor/internal/infra/storage"
	"linux-tutor/internal/terminal"
	"os"
)

func Run() {
	db, err := storage.Open("data/linux-tutor.db")
	if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
	defer db.Close()
	if err := storage.Migrate(db); err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
	if err := terminal.Start(db); err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
}
