package main

import (
	"log"
	"github.com/nihar-hegde/valtro-backend/internal/database"
	"github.com/nihar-hegde/valtro-backend/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: .env file not found, reading from system environment")
	}

	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create and start the server
	s := server.NewServer(db)
	if err := s.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}