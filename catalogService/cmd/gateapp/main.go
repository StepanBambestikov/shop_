package main

import (
	"catalogServiceGit/cmd/gateapp/cmd"
	"log"
)

func main() {
	log.Println("Starting counter gate")
	if err := cmd.Execute(); err != nil {
		log.Fatal("Exiting due to error: %w", err)
	}
}
