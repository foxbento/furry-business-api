package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/foxbento/furry-business-api/config"
	"github.com/foxbento/furry-business-api/db"
	"github.com/foxbento/furry-business-api/handlers"
	"github.com/rs/cors"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := db.Initialize(cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/businesses", handlers.GetBusinesses).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Adjust this in production
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, c.Handler(r)))
}