package main

import (
	"catalogServiceGit/cmd/coreapp/cmd"
	"log"
)

func main() {
	log.Println("Starting counter core")
	if err := cmd.Execute(); err != nil {
		log.Fatal("Exiting due to error: %w", err)
	}
}
