package main

import "linux-tutor/internal/app"

func main() {
	if err := app.RunTUI(); err != nil {
		panic(err)
	}
}
