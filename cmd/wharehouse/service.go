package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"wbtechschool-L3/WarehouseControl/internal/config"
	"wbtechschool-L3/WarehouseControl/internal/handler"
	appmiddleware "wbtechschool-L3/WarehouseControl/internal/middleware"
	"wbtechschool-L3/WarehouseControl/internal/repository"
	"wbtechschool-L3/WarehouseControl/internal/service"
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

	// Initialize service layer
	authService := service.NewAuthService(repo, cfg.JWT)
	itemService := service.NewItemService(repo)
	historyService := service.NewHistoryService(repo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	itemHandler := handler.NewItemHandler(itemService, authService)
	historyHandler := handler.NewHistoryHandler(historyService, authService)

	// Setup router with middleware
	r := chi.NewRouter()

	// Global middleware
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Timeout(60 * time.Second))
	r.Use(appmiddleware.CORSMiddleware)

	// Public routes (no authentication required)
	r.Post("/auth/login", authHandler.Login)

	// Serve static files
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})

	// Protected routes (authentication required)
	r.Group(func(r chi.Router) {
		r.Use(appmiddleware.AuthMiddleware(authService))

		// Item routes
		r.Route("/items", func(r chi.Router) {
			r.Get("/", itemHandler.GetAllItems)      // GET /items
			r.Post("/", itemHandler.CreateItem)       // POST /items
			r.Get("/{id}", itemHandler.GetItem)       // GET /items/{id}
			r.Put("/{id}", itemHandler.UpdateItem)    // PUT /items/{id}
			r.Delete("/{id}", itemHandler.DeleteItem) // DELETE /items/{id}

			// History routes
			r.Get("/{id}/history", historyHandler.GetHistoryByItemID) // GET /items/{id}/history
		})

		// History routes
		r.Get("/history", historyHandler.GetAllHistory) // GET /history
	})

	// Start HTTP server
	log.Printf("Starting HTTP server on %s", cfg.HTTPServer.Address)
	if err := http.ListenAndServe(cfg.HTTPServer.Address, r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

	return nil
}
