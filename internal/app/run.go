package app

import (
	"fmt"
	"linux-tutor/internal/terminal"
	"os"
)

func Run() {
	if err := terminal.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
