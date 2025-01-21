package api

import (
	"github.com/EgSundqvist/imagebook-imageapi/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Lägg till CORS-middleware
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
		protected.POST("/upload-url", GenerateUploadURLHandler)         // Endpoint för att generera presignerade URL:er för uppladdning
		protected.POST("/confirm-upload", ConfirmUploadHandler)         // Endpoint för att bekräfta uppladdning och spara metadata
		protected.POST("/delete-url", GenerateDeleteURLHandler)         // Endpoint för att generera presignerade URL:er för borttagning
		protected.POST("/confirm-delete", ConfirmDeleteHandler)         // Endpoint för att bekräfta borttagning och ta bort metadata
		protected.GET("/images/page/:pageNumber", GetImagesPageHandler) // Skyddad endpoint för paginerade bilder
		protected.GET("/images/:id", GetImageByIdHandler)

	}
	router.POST("/create-s3-folder", services.BackendAuthMiddleware(), CreateS3FolderHandler)

	return router
}
