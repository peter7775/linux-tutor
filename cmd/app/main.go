package main

import "linux-tutor/internal/app"

func main() {
	err := app.RunGUI()
	if err != nil {
		return
	}
}
