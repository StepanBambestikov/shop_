package main

import (
	"log"
	"orderServiceGit/cmd/gateapp/cmd"
)

func main() {
	log.Println("Starting counter gate")
	if err := cmd.Execute(); err != nil {
		log.Fatal("Exiting due to error: %w", err)
	}
}
