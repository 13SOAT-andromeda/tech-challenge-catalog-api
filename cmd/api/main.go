package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	awsadapter "github.com/juliovaz/tech-challenge-catalog-api/internal/adapter/aws"
	maintenancemodel "github.com/juliovaz/tech-challenge-catalog-api/internal/adapter/database/model/maintenance"
	productmodel "github.com/juliovaz/tech-challenge-catalog-api/internal/adapter/database/model/product"
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
	if err := db.AutoMigrate(&productmodel.Product{}, &maintenancemodel.Maintenance{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// AWS Setup
	snsPublisher, err := awsadapter.NewSNSPublisher(context.Background())
	if err != nil {
		log.Fatalf("failed to setup SNS publisher: %v", err)
	}

	// Wire: repo -> service -> handler
	productRepo := repository.NewProductRepository(db)
	productSvc := services.NewProductService(productRepo, snsPublisher)
	productH := handlers.NewProductHandler(productSvc)

	maintenanceRepo := repository.NewMaintenanceRepository(db)
	maintenanceSvc := services.NewMaintenanceService(maintenanceRepo)
	maintenanceH := handlers.NewMaintenanceHandler(maintenanceSvc)

	router := adapthttp.SetupRouter(productH, maintenanceH)

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
