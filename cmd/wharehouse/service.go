package main

import (
	"log"
	"wbtechschool-L3/WarehouseControl/internal/config"
)

// Run initializes and starts the warehouse application
func Run(configPath string) error {
	// Load configuration from config file and environment variables
	cfg := config.MustLoad(configPath)

	log.Printf("Starting warehouse application with config:")
	log.Printf("  Environment: %s", cfg.Env)
	log.Printf("  Server Address: %s", cfg.HTTPServer.Address)
	log.Printf("  Database: %s@%s:%d/%s", cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

	// TODO: Initialize database connection
	// TODO: Initialize repository layer
	// TODO: Initialize service layer
	// TODO: Initialize handlers
	// TODO: Setup router with middleware
	// TODO: Start HTTP server

	return nil
}
