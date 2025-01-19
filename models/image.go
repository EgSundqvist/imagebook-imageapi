package models

import "gorm.io/gorm"

type Image struct {
	gorm.Model
	URL          string `gorm:"not null"`
	Description  string
	UserID       *uint
	PresignedURL string `gorm:"-"`
}
