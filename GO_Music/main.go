package main

import (
	"log"

	"GO_Music/app"
)

func main() {
	if err := app.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}