package api

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/EgSundqvist/imagebook-imageapi/data"
	"github.com/EgSundqvist/imagebook-imageapi/services"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

// ...existing code...

func GetImagesPageHandler(c *gin.Context) {
	pageNumber, err := strconv.Atoi(c.Param("pageNumber"))
	if err != nil || pageNumber < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

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

	const pageSize = 3
	images, err := data.GetImagesPage(userID, pageNumber, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch images"})
		return
	}

	sess, err := services.InitAWSSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create AWS session"})
		return
	}

	svc := s3.New(sess)
	bucket := "imagebook-images"

	for i := range images {
		// Parse the URL to extract the key
		parsedURL, err := url.Parse(images[i].URL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse image URL"})
			return
		}
		key := strings.TrimPrefix(parsedURL.Path, "/")

		// Generate a presigned URL
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
		images[i].PresignedURL = presignedURL
	}

	c.JSON(http.StatusOK, images)
}

// ...existing code...
