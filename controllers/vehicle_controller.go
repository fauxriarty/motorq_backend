package controllers

import (
	"motorq_backend/database"
	"motorq_backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateVehicle(c *gin.Context) {
	var input models.Vehicle
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vehicle := models.Vehicle{
		Make:         input.Make,
		Model:        input.Model,
		LicensePlate: input.LicensePlate,
		Available:    true,
	}

	database.DB.Create(&vehicle)

	c.JSON(http.StatusOK, gin.H{"data": vehicle})
}

func GetVehicles(c *gin.Context) {
	var vehicles []models.Vehicle
	database.DB.Find(&vehicles)

	c.JSON(http.StatusOK, gin.H{"data": vehicles})
}
