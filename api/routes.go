package api

import (
	"github.com/EgSundqvist/imagebook-imageapi/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5000", "https://userapi.spacetechnology.net", "https://imagebook.spacetechnology.net"}, // Tillåt anrop från dessa URL:er (lägg till user-API:ets URL)
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Lägg till en enkel ping endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Skyddade endpoints
	protected := router.Group("/")
	protected.Use(services.AuthMiddleware())
	{
		protected.POST("/upload-url", GenerateUploadURLHandler)
		protected.POST("/confirm-upload", ConfirmUploadHandler)
		protected.POST("/delete-url", GenerateDeleteURLHandler)
		protected.POST("/confirm-delete", ConfirmDeleteHandler)
		protected.GET("/images/page/:pageNumber", GetImagesPageHandler)
		protected.GET("/images/:id", GetImageByIdHandler)

	}
	router.POST("/create-s3-folder", services.BackendAuthMiddleware(), CreateS3FolderHandler)

	return router
}
