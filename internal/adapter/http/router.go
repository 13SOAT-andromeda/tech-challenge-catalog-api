package http

import (
	"github.com/gin-gonic/gin"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/adapter/http/handlers"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/adapter/http/middlewares"
)

func SetupRouter(categoryHandler *handlers.CategoryHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	{
		categories := api.Group("/categories")
		{
			categories.POST("/", categoryHandler.Create)
			categories.GET("/", categoryHandler.List)
			categories.GET("/:id", categoryHandler.GetByID)

			protected := categories.Group("/")
			protected.Use(middlewares.AuthRequired())
			{
				protected.PUT("/:id", categoryHandler.Update)
				protected.DELETE("/:id", categoryHandler.Delete)
			}
		}
	}
	return r
}
