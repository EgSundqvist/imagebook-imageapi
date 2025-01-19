package data

import (
	"log"

	"github.com/EgSundqvist/imagebook-imageapi/models"
)

func Seed() {
	// Check if there are already images in the database
	var count int64
	if err := DB.Model(&models.Image{}).Count(&count).Error; err != nil {
		log.Fatalf("Failed to count images: %v", err)
	}

	// If there are no images, seed the database
	if count == 0 {
		images := []models.Image{
			{URL: "https://imagebook-images.s3.eu-north-1.amazonaws.com/users/defaultUser/images/20150103_124830.jpg", Description: "Image 1"},
			{URL: "https://imagebook-images.s3.eu-north-1.amazonaws.com/users/defaultUser/images/20150103_145629.jpg", Description: "Image 2"},
			{URL: "https://imagebook-images.s3.eu-north-1.amazonaws.com/users/defaultUser/images/20150103_151620.jpg", Description: "Image 3"},
			{URL: "https://imagebook-images.s3.eu-north-1.amazonaws.com/users/defaultUser/images/20150107_020628.jpg", Description: "Image 4"},
			{URL: "https://imagebook-images.s3.eu-north-1.amazonaws.com/users/defaultUser/images/20150107_025553.jpg", Description: "Image 5"},
			{URL: "https://imagebook-images.s3.eu-north-1.amazonaws.com/users/defaultUser/images/20150107_151824.jpg", Description: "Image 6"},
			{URL: "https://imagebook-images.s3.eu-north-1.amazonaws.com/users/defaultUser/images/IMG_20141225_103802.jpg", Description: "Image 7"},
		}

		for _, image := range images {
			DB.Create(&image)
		}
	}
}
