package api

import (
	"log"
	"net/http"

	"github.com/EgSundqvist/imagebook-imageapi/config"
	"github.com/EgSundqvist/imagebook-imageapi/services"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

func CreateS3FolderHandler(c *gin.Context) {

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
	bucket := config.AppConfig.S3Bucket
	userFolderKey := "users/" + userID.(string) + "/"
	imagesFolderKey := userFolderKey + "images/"

	// Skapa "user"-mappen i S3
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(userFolderKey),
	})
	if err != nil {
		log.Printf("Failed to create user folder in S3: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user folder in S3"})
		return
	}

	// Skapa "images"-mappen i S3
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
