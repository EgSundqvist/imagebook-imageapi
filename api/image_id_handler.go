package api

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/EgSundqvist/imagebook-imageapi/config"
	"github.com/EgSundqvist/imagebook-imageapi/data"
	"github.com/EgSundqvist/imagebook-imageapi/services"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

func GetImageByIdHandler(c *gin.Context) {
	id := c.Param("id")

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

	image, err := data.GetImageByID(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	sess, err := services.InitAWSSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create AWS session"})
		return
	}

	svc := s3.New(sess)
	bucket := config.AppConfig.S3Bucket

	parsedURL, err := url.Parse(image.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse image URL"})
		return
	}
	key := strings.TrimPrefix(parsedURL.Path, "/")

	// Generera en presignerad URL för att hämta objektet från S3
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	presignedURL, err := req.Presign(15 * time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate presigned URL"})
		return
	}

	// Inkludera presignedURL i image-objektet
	image.PresignedURL = presignedURL

	c.JSON(http.StatusOK, image)
}
