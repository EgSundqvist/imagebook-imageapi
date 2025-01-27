package data

import (
	"log"

	"github.com/EgSundqvist/imagebook-imageapi/models"
)

func Seed() {
	var count int64
	if err := DB.Model(&models.Image{}).Count(&count).Error; err != nil {
		log.Fatalf("Failed to count images: %v", err)
	}

	if count == 0 {
		images := []models.Image{
			{URL: "your-seed-image1.jpg", Description: "Image 1"},
			{URL: "your-seed-image2.jpg", Description: "Image 2"},
			{URL: "your-seed-image3.jpg", Description: "Image 3"},
			{URL: "your-seed-image4.jpg", Description: "Image 4"},
			{URL: "your-seed-image5.jpg", Description: "Image 5"},
			{URL: "your-seed-image6.jpg", Description: "Image 6"},
			{URL: "your-seed-image7.jpg", Description: "Image 7"},
		}

		for _, image := range images {
			DB.Create(&image)
		}
	}
}
