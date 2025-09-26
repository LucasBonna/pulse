package main

import (
	"log"

	"lucasbonna/pulse/internal/api"
	"lucasbonna/pulse/internal/storage"
)

func main() {
	// Initialize database
	db, err := storage.NewSQLiteDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize HTTP server with dependencies
	server := api.NewServer(db)

	// Start server
	log.Println("Starting server on :8080...")
	if err := server.Start(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
