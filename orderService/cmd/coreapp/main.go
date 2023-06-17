package main

import (
	"log"
	"orderServiceGit/cmd/coreapp/cmd"
)

func main() {
	log.Println("Starting counter core")
	if err := cmd.Execute(); err != nil {
		log.Fatal("Exiting due to error: %w", err)
	}
}
