package main

import (
	"log"
	"net/http"

	"github.com/Ruseigha/sendlink/internal/cache"
	"github.com/Ruseigha/sendlink/internal/config"
	"github.com/Ruseigha/sendlink/internal/database"
	"github.com/Ruseigha/sendlink/internal/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
    // Load .env file if it exists (for local development)
    godotenv.Load()

    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    log.Printf("Starting ShortLink API in %s environment", cfg.Environment)

    // Initialize database
    db, err := database.NewPostgresDB(cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    if err := db.InitSchema(); err != nil {
        log.Fatalf("Failed to initialize schema: %v", err)
    }

    // Initialize cache
    redisCache, err := cache.NewRedisCache(cfg.RedisURL, cfg.RedisPassword)
    if err != nil {
        log.Printf("Warning: Redis not available: %v", err)
        // Continue without cache
    }

    // Initialize handlers
    h := handlers.NewHandler(db, redisCache, cfg)

    // Setup router
    router := mux.NewRouter()
    router.HandleFunc("/health", h.HealthCheck).Methods("GET")
    router.HandleFunc("/api/shorten", h.CreateShortURL).Methods("POST")
    router.HandleFunc("/api/stats", h.GetStats).Methods("GET")
    router.HandleFunc("/{code}", h.RedirectURL).Methods("GET")

    // Start server
    log.Printf("Server listening on port %s", cfg.ServerPort)
    if err := http.ListenAndServe(":"+cfg.ServerPort, router); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}