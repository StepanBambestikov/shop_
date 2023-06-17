package main

import (
	"gitea.teneshag.ru/gigabit/goauth/cmd/goauth/cmd"
	"log"
)

func main() {
	log.Println("Starting goauth")
	if err := cmd.Execute(); err != nil {
		log.Fatal("Exiting due to error: %w", err)
	}
}
