package main

import (
	"log"
	"wbtechschool-L3/WarehouseControl/internal/config"
	"wbtechschool-L3/WarehouseControl/internal/repository"
)

// Run initializes and starts the warehouse application
func Run(configPath string) error {
	// Load configuration from config file and environment variables
	cfg := config.MustLoad(configPath)

	log.Printf("Starting warehouse application with config:")
	log.Printf("  Environment: %s", cfg.Env)
	log.Printf("  Server Address: %s", cfg.HTTPServer.Address)
	log.Printf("  Database: %s@%s:%d/%s", cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

	// Initialize repository layer with database connection
	repo, err := repository.New(cfg.Database)
	if err != nil {
		log.Fatalf("failed to initialize repository: %v", err)
	}
	defer repo.Close()

	log.Println("Database connection established successfully")

	// TODO: Initialize service layer
	// TODO: Initialize handlers
	// TODO: Setup router with middleware
	// TODO: Start HTTP server

	return nil
}
