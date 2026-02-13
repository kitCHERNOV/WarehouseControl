package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: failed to load .env file: %v", err)
	}

	// Get config path from environment variable, with fallback to default
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/local.yaml"
	}

	// Allow command line override
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	// Run the warehouse application
	if err := Run(configPath); err != nil {
		log.Fatalf("failed to run application: %v", err)
	}
}
