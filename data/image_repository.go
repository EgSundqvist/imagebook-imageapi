package data

import (
	"github.com/EgSundqvist/imagebook-imageapi/models"
)

func GetAllImages() ([]models.Image, error) {
	var images []models.Image
	if err := DB.Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

func GetImageByID(id string, userID uint) (models.Image, error) {
	var image models.Image
	if err := DB.Where("id = ? AND user_id = ?", id, userID).First(&image).Error; err != nil {
		return image, err
	}
	return image, nil
}

func GetImagesPage(userID uint, pageNumber int, pageSize int) ([]models.Image, error) {
	var images []models.Image
	offset := (pageNumber - 1) * pageSize
	if err := DB.Where("user_id = ?", userID).Order("id").Offset(offset).Limit(pageSize).Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

func CreateImage(image models.Image) error {
	if err := DB.Create(&image).Error; err != nil {
		return err
	}
	return nil
}

func DeleteImageByID(id string, userID uint) error {
	if err := DB.Unscoped().Where("id = ? AND user_id = ?", id, userID).Delete(&models.Image{}).Error; err != nil {
		return err
	}
	return nil
}
