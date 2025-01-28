package api

import (
	"os"

	"github.com/EgSundqvist/imagebook-imageapi/config"
	"github.com/EgSundqvist/imagebook-imageapi/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	var allowOrigins []string
	if os.Getenv("RUNENVIRONMENT") == "Production" {
		allowOrigins = []string{config.AppConfig.UserAPIURL, config.AppConfig.ClientURL}
	} else {
		allowOrigins = []string{os.Getenv("USERAPI_URL"), os.Getenv("CLIENT_URL")}
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

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
