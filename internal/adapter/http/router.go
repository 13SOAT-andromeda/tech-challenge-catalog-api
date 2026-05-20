package http

import (
	"github.com/gin-gonic/gin"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/adapter/http/handlers"
)

func SetupRouter(productH *handlers.ProductHandler, maintenanceH *handlers.MaintenanceHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
  
	api := r.Group("/")
	{
		products := api.Group("/products")
		{
			products.POST("/", productH.Create)
			products.GET("/", productH.List)
			products.GET("/:id", productH.GetByID)
			products.PUT("/:id", productH.Update)
			products.DELETE("/:id", productH.Delete)
		}

		maintenances := api.Group("/maintenances")
		{
			maintenances.POST("/", maintenanceH.Create)
			maintenances.GET("/", maintenanceH.List)
			maintenances.GET("/:id", maintenanceH.GetByID)
			maintenances.PUT("/:id", maintenanceH.Update)
			maintenances.DELETE("/:id", maintenanceH.Delete)
		}
	}

	return r
}
