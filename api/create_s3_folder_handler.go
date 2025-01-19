package api

import (
	"log"
	"net/http"

	"github.com/EgSundqvist/imagebook-imageapi/services"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

// CreateS3FolderRequest representerar förfrågan för att skapa en ny mapp i S3
type CreateS3FolderRequest struct {
}

// CreateS3FolderHandler skapar en ny mapp i S3 för en användare
// ...existing code...

func CreateS3FolderHandler(c *gin.Context) {
	// Ta bort bindningen av JSON-request eftersom vi inte längre behöver UserID från requesten
	/*
	   var req CreateS3FolderRequest
	   if err := c.ShouldBindJSON(&req); err != nil {
	       log.Printf("Invalid request: %v", err)
	       c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
	       return
	   }
	*/

	// Hämta userID från JWT-token
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	log.Printf("Received request to create S3 folder for user ID: %s", userID)

	sess, err := services.InitAWSSession()
	if err != nil {
		log.Printf("Failed to create AWS session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create AWS session"})
		return
	}

	svc := s3.New(sess)
	bucket := "imagebook-images"
	userFolderKey := "users/" + userID.(string) + "/"
	imagesFolderKey := userFolderKey + "images/"

	// Create the user folder in S3
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(userFolderKey),
	})
	if err != nil {
		log.Printf("Failed to create user folder in S3: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user folder in S3"})
		return
	}

	// Create the images folder in S3
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(imagesFolderKey),
	})
	if err != nil {
		log.Printf("Failed to create images folder in S3: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create images folder in S3"})
		return
	}

	log.Printf("Created S3 folders for user ID: %s", userID)

	c.JSON(http.StatusCreated, gin.H{"message": "S3 folders created successfully"})
}

// ...existing code...
