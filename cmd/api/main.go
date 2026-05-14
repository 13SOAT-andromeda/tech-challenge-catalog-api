package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	categorymodel "github.com/juliovaz/tech-challenge-catalog-api/internal/adapter/database/model/category"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/adapter/database/repository"
	adapthttp "github.com/juliovaz/tech-challenge-catalog-api/internal/adapter/http"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/adapter/http/handlers"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/application/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Auto Migration
	if err := db.AutoMigrate(&categorymodel.Category{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// Wire Category: repo -> service -> handler
	categoryRepo := repository.NewCategoryRepository(db)
	categorySvc := services.NewCategoryService(categoryRepo)
	categoryH := handlers.NewCategoryHandler(categorySvc)

	router := adapthttp.SetupRouter(categoryH)

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
