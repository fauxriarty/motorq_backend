package main

import (
	"motorq_backend/database"
	"motorq_backend/models"
	"motorq_backend/routes"
)

func main() {
	database.ConnectDatabase()
	database.DB.AutoMigrate(&models.Driver{}, &models.Vehicle{}, &models.Assignment{})

	r := routes.SetupRouter()

	r.Run()
}
