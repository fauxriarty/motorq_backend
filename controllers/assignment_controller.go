package controllers

import (
	"motorq_backend/database"
	"motorq_backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AssignVehicleToDriver(c *gin.Context) {
	var input models.Assignment
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var driver models.Driver
	var vehicle models.Vehicle

	if err := database.DB.First(&driver, input.DriverID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Driver not found!"})
		return
	}

	if err := database.DB.First(&vehicle, input.VehicleID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vehicle not found!"})
		return
	}

	if !vehicle.Available {
		c.JSON(http.StatusConflict, gin.H{"error": "Vehicle is already assigned!"})
		return
	}

	vehicle.Available = false
	database.DB.Save(&vehicle)

	assignment := models.Assignment{
		DriverID:  input.DriverID,
		VehicleID: input.VehicleID,
		StartTime: time.Now(),
		Status:    "Active",
	}

	database.DB.Create(&assignment)

	c.JSON(http.StatusOK, gin.H{"data": assignment})
}

func UnassignVehicleFromDriver(c *gin.Context) {
	var input models.Assignment
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var assignment models.Assignment
	if err := database.DB.Where("driver_id = ? AND vehicle_id = ? AND status = 'Active'", input.DriverID, input.VehicleID).First(&assignment).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Assignment not found!"})
		return
	}

	// Update vehicle availability
	var vehicle models.Vehicle
	database.DB.First(&vehicle, assignment.VehicleID)
	vehicle.Available = true
	database.DB.Save(&vehicle)

	// Mark assignment as completed
	assignment.Status = "Completed"
	assignment.EndTime = time.Now()
	database.DB.Save(&assignment)

	c.JSON(http.StatusOK, gin.H{"data": assignment})
}
