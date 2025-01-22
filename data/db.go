package data

import (
	"fmt"
	"log"

	"github.com/EgSundqvist/imagebook-imageapi/models"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase(file, server, database, username, password string, port int) {
	if len(file) == 0 {
		DB = openMySql(server, database, username, password, port)
	} else {
		DB, _ = gorm.Open(sqlite.Open(file), &gorm.Config{})
	}
	DB.AutoMigrate(&models.Image{})
}

func openMySql(server, database, username, password string, port int) *gorm.DB {
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, server, port, database)

	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return db
}
