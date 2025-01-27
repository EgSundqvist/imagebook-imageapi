package main

import (
	"github.com/EgSundqvist/imagebook-imageapi/api"
	"github.com/EgSundqvist/imagebook-imageapi/config"
	"github.com/EgSundqvist/imagebook-imageapi/data"
)

func main() {
	config.LoadConfig()

	data.InitDatabase(
		config.AppConfig.Database.File,
		config.AppConfig.Database.Server,
		config.AppConfig.Database.Database,
		config.AppConfig.Database.Username,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.Port,
	)

	router := api.SetupRouter()
	router.Run(":8080")
}
