package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/Bughay/Trainer-GO/db"
	"github.com/Bughay/Trainer-GO/internal/auth"
	"github.com/Bughay/Trainer-GO/internal/food"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it:", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	ctx := context.Background()

	dbPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal("Failed to create connection pool:", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	queries := db.New(dbPool)

	authHandler, err := auth.NewAuthHandler(queries, jwtSecret)
	foodHandler := food.NewFoodHandler(queries)
	if err != nil {
		log.Fatalf("Failed to create auth handler: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/auth/register", authHandler.UserRegistrationHandler)
	mux.HandleFunc("/auth/login", authHandler.UserLoginHandler)
	mux.HandleFunc("/food/create", authHandler.AuthMiddleware(foodHandler.CreateFoodItemHandler))
	mux.HandleFunc("/food/log", authHandler.AuthMiddleware(foodHandler.LogFoodHandler))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Println("Server starting on :8080...")
	log.Fatal(server.ListenAndServe())
}
