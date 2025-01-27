package api

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/EgSundqvist/imagebook-imageapi/config"
	"github.com/EgSundqvist/imagebook-imageapi/data"
	"github.com/EgSundqvist/imagebook-imageapi/services"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

type DeleteURLRequest struct {
	ID string `json:"id" binding:"required"`
}

func GenerateDeleteURLHandler(c *gin.Context) {
	var req DeleteURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	log.Printf("Received request to generate delete URL for ID: %s", req.ID)

	userIDStr := c.GetString("userID")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userIDUint64, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user ID"})
		return
	}
	userID := uint(userIDUint64)

	image, err := data.GetImageByID(req.ID, userID)
	if err != nil {
		log.Printf("Failed to get image by ID: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	sess, err := services.InitAWSSession()
	if err != nil {
		log.Printf("Failed to create AWS session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create AWS session"})
		return
	}

	svc := s3.New(sess)
	bucket := config.AppConfig.S3Bucket
	key := image.URL[len("https://"+bucket+".s3.eu-north-1.amazonaws.com/"):]

	// Generate a presigned URL for DELETE
	delReq, _ := svc.DeleteObjectRequest(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	presignedURL, err := delReq.Presign(15 * time.Minute)
	if err != nil {
		log.Printf("Failed to generate presigned URL: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate presigned URL"})
		return
	}

	log.Printf("Generated presigned URL: %s", presignedURL)

	// Returnera den presignerade URL:en till klienten
	c.JSON(http.StatusOK, gin.H{"PresignedURL": presignedURL})
}

func ConfirmDeleteHandler(c *gin.Context) {
	var req struct {
		ID string `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	log.Printf("Received confirmation for delete: ID: %s", req.ID)

	userIDStr := c.GetString("userID")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userIDUint64, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user ID"})
		return
	}
	userID := uint(userIDUint64)

	// Kontrollera att bilden existerar och tillhör användaren
	if _, err := data.GetImageByID(req.ID, userID); err != nil {
		log.Printf("Failed to get image by ID: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	// Ta bort bildens metadata från databasen
	if err := data.DeleteImageByID(req.ID, userID); err != nil {
		log.Printf("Failed to delete image metadata: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image metadata"})
		return
	}

	log.Printf("Image metadata deleted successfully")

	c.JSON(http.StatusOK, gin.H{"message": "Image metadata deleted successfully"})
}
