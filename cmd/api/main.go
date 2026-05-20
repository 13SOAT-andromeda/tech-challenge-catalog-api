package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	awsadapter "github.com/juliovaz/tech-challenge-catalog-api/internal/adapter/aws"
	backordermodel "github.com/juliovaz/tech-challenge-catalog-api/internal/adapter/database/model/backorder"
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

	if err := db.AutoMigrate(
		&productmodel.Product{},
		&maintenancemodel.Maintenance{},
		&backordermodel.Backorder{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	snsPublisher, err := awsadapter.NewSNSPublisher(ctx)
	if err != nil {
		log.Fatalf("failed to setup SNS publisher: %v", err)
	}

	var emailPublisher *awsadapter.EmailNotificationPublisher
	if ep, err := awsadapter.NewEmailNotificationPublisher(ctx); err != nil {
		log.Printf("email notification publisher not configured: %v", err)
	} else {
		emailPublisher = ep
	}

	productRepo := repository.NewProductRepository(db)
	backorderRepo := repository.NewBackorderRepository(db)
	productSvc := services.NewProductService(productRepo, backorderRepo, snsPublisher, emailPublisher)

	sqsConsumer, err := awsadapter.NewSQSConsumer(ctx, productSvc)
	if err != nil {
		log.Printf("SQS consumer not started: %v", err)
	} else {
		go sqsConsumer.Start(ctx)
	}

	maintenanceRepo := repository.NewMaintenanceRepository(db)
	maintenanceSvc := services.NewMaintenanceService(maintenanceRepo)
	maintenanceH := handlers.NewMaintenanceHandler(maintenanceSvc)

	productH := handlers.NewProductHandler(productSvc)
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
