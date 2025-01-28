package api

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/EgSundqvist/imagebook-imageapi/config"
	"github.com/EgSundqvist/imagebook-imageapi/data"
	"github.com/EgSundqvist/imagebook-imageapi/models"
	"github.com/EgSundqvist/imagebook-imageapi/services"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

type UploadURLRequest struct {
	Filename    string `json:"filename" binding:"required"`
	Description string `json:"description"`
}

func GenerateUploadURLHandler(c *gin.Context) {
	var req UploadURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		log.Printf("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	log.Printf("Received request to generate upload URL for user ID: %s, filename: %s, description: %s", userID, req.Filename, req.Description)

	sess, err := services.InitAWSSession()
	if err != nil {
		log.Printf("Failed to create AWS session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create AWS session"})
		return
	}

	svc := s3.New(sess)
	bucket := config.AppConfig.S3Bucket
	key := "users/" + userID + "/images/" + req.Filename

	putReq, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	presignedURL, err := putReq.Presign(15 * time.Minute)
	if err != nil {
		log.Printf("Failed to generate presigned URL: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate presigned URL"})
		return
	}

	log.Printf("Generated presigned URL: %s", presignedURL)

	// Returnera den presignerade URL:en till klienten
	c.JSON(http.StatusOK, gin.H{"presignedURL": presignedURL, "s3URL": "https://" + bucket + ".s3.eu-north-1.amazonaws.com/" + key, "description": req.Description})
}

// ConfirmUploadHandler bekräftar uppladdningen och sparar metadata i databasen
func ConfirmUploadHandler(c *gin.Context) {
	var req struct {
		S3URL       string `json:"s3URL" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userIDStr := c.GetString("userID")
	if userIDStr == "" {
		log.Printf("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userIDUint64, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		log.Printf("Failed to parse user ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user ID"})
		return
	}
	userID := uint(userIDUint64)

	log.Printf("Received confirmation for upload: user ID: %s, S3URL: %s, Description: %s", userIDStr, req.S3URL, req.Description)

	// Spara bildens metadata i databasen
	image := models.Image{
		URL:         req.S3URL,
		Description: req.Description,
		UserID:      &userID,
	}
	if err := data.CreateImage(image); err != nil {
		log.Printf("Failed to save image metadata: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image metadata"})
		return
	}

	log.Printf("Image metadata saved successfully")

	c.JSON(http.StatusOK, gin.H{"message": "Image metadata saved successfully"})
}
