package controllers

import (
	"log"
	"motorq_backend/database"
	"motorq_backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AssignVehicleToDriver(c *gin.Context) {
	var request struct {
		DriverID  uint   `json:"driver_id"`
		VehicleID uint   `json:"vehicle_id"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var driver models.Driver
	if err := database.DB.First(&driver, request.DriverID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
		return
	}

	var vehicle models.Vehicle
	if err := database.DB.First(&vehicle, request.VehicleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
		return
	}

	startTime, err := time.Parse("15:04", request.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time format"})
		return
	}

	endTime, err := time.Parse("15:04", request.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end time format"})
		return
	}

	// Ensure no conflicts with existing assignments
	var conflictCount int64
	database.DB.Model(&models.Assignment{}).Where("driver_id = ? AND ((start_time < ? AND end_time > ?) OR (start_time < ? AND end_time > ?))",
		request.DriverID, endTime, startTime, endTime, startTime).Count(&conflictCount)

	if conflictCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Driver is already assigned to another vehicle during this time period"})
		return
	}

	// Proceed with assignment
	assignment := models.Assignment{
		DriverID:  request.DriverID,
		VehicleID: request.VehicleID,
		StartTime: startTime,
		EndTime:   endTime,
	}

	if err := database.DB.Create(&assignment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create assignment"})
		return
	}

	// Update the vehicle's AssignedDriverID field
	vehicle.AssignedDriverID = &driver.ID
	if err := database.DB.Save(&vehicle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vehicle's assigned driver"})
		return
	}

	// Update driver status if needed
	if startTime.Before(time.Now()) && endTime.After(time.Now()) {
		driver.Status = "busy"
	} else {
		driver.Status = "available"
	}

	if err := database.DB.Save(&driver).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update driver status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": assignment})
}

func UnassignVehicleFromDriver(c *gin.Context) {
	var input struct {
		VehicleID uint `json:"vehicle_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var vehicle models.Vehicle
	if err := database.DB.First(&vehicle, input.VehicleID).Error; err != nil {
		log.Println("Vehicle not found:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
		return
	}

	if vehicle.AssignedDriverID == nil {
		log.Println("No driver assigned to this vehicle")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No driver assigned to this vehicle"})
		return
	}

	// Find the assignment related to this vehicle and driver, and ensure the time slot matches
	var assignment models.Assignment
	if err := database.DB.Where("vehicle_id = ? AND driver_id = ?", vehicle.ID, *vehicle.AssignedDriverID).First(&assignment).Error; err != nil {
		log.Println("Assignment not found:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		return
	}

	log.Println("Assignment Data:", assignment)

	// Perform the deletion
	if err := database.DB.Delete(&assignment).Error; err != nil {
		log.Println("Failed to delete assignment:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete assignment"})
		return
	}

	log.Println("Assignment deleted successfully")

	// Update the driver status to "available"
	var driver models.Driver
	if err := database.DB.First(&driver, *vehicle.AssignedDriverID).Error; err == nil {
		driver.Status = "available"
		if err := database.DB.Save(&driver).Error; err != nil {
			log.Println("Failed to update driver status:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update driver status"})
			return
		}
	}

	// Unassign the driver from the vehicle
	vehicle.AssignedDriverID = nil
	if err := database.DB.Save(&vehicle).Error; err != nil {
		log.Println("Failed to unassign driver from vehicle:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unassign driver from vehicle"})
		return
	}

	log.Println("Driver unassigned from vehicle and assignment deleted successfully")
	c.JSON(http.StatusOK, gin.H{"data": "Driver unassigned from vehicle and assignment deleted successfully"})
}
