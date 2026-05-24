package main

import "linux-tutor/internal/app"

func main() {
	if err := app.RunGUI(); err != nil {
		panic(err)
	}
}
