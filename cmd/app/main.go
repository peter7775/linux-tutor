package main

import (
	"linux-tutor/internal/app"
	"log"
)

func main() {
	if err := app.RunGUI(); err != nil {
		log.Fatal(err)
	}
}
