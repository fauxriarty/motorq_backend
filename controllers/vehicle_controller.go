package controllers

import (
	"log"
	"motorq_backend/database"
	"motorq_backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateVehicle(c *gin.Context) {
	var input models.Vehicle
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error binding JSON:", err) // Log the error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vehicle := models.Vehicle{
		Make:             input.Make,
		Model:            input.Model,
		LicensePlate:     input.LicensePlate,
		AssignedDriverID: input.AssignedDriverID,
	}

	if err := database.DB.Create(&vehicle).Error; err != nil {
		log.Println("Error creating vehicle:", err) // Log the error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vehicle"})
		return
	}

	log.Println("Vehicle created successfully:", vehicle) // Log successful creation
	c.JSON(http.StatusOK, gin.H{"data": vehicle})
}

func GetVehicles(c *gin.Context) {
	var vehicles []models.Vehicle
	result := database.DB.Find(&vehicles)

	if result.Error != nil {
		log.Println("Error retrieving vehicles:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve vehicles"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"data": []models.Vehicle{}})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": vehicles})
	}
}
