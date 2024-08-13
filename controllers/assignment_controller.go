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

	// Check if the requested time slot is within the driver's working hours
	if startTime.Hour() < 9 || endTime.Hour() > 17 {
		c.JSON(http.StatusConflict, gin.H{"error": "Assignment time is outside of driver's working hours"})
		return
	}

	// Ensure no conflicts with existing vehicle assignments
	var vehicleConflictCount int64
	database.DB.Model(&models.Assignment{}).Where("vehicle_id = ? AND ((start_time < ? AND end_time > ?) OR (start_time < ? AND end_time > ?))",
		request.VehicleID, endTime, startTime, endTime, startTime).Count(&vehicleConflictCount)

	if vehicleConflictCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "This vehicle is already booked during this time period"})
		return
	}

	// Proceed with assignment
	assignment := models.Assignment{
		DriverID:  request.DriverID,
		VehicleID: request.VehicleID,
		StartTime: startTime,
		EndTime:   endTime,
		Status:    "pending",
	}

	if err := database.DB.Create(&assignment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create assignment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": assignment})
}

func UnassignVehicleFromDriver(c *gin.Context) {
	var input struct {
		DriverID  uint `json:"driver_id"`
		VehicleID uint `json:"vehicle_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var assignment models.Assignment
	if err := database.DB.Where("vehicle_id = ? AND driver_id = ?", input.VehicleID, input.DriverID).First(&assignment).Error; err != nil {
		log.Println("Assignment not found:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		return
	}

	if err := database.DB.Delete(&assignment).Error; err != nil {
		log.Println("Failed to delete assignment:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete assignment"})
		return
	}

	log.Println("Assignment deleted successfully")
	c.JSON(http.StatusOK, gin.H{"data": "Driver unassigned from vehicle and assignment deleted successfully"})
}

func AcceptAssignment(c *gin.Context) {
	var input struct {
		DriverID     uint `json:"driver_id"`
		AssignmentID uint `json:"assignment_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var assignment models.Assignment
	if err := database.DB.First(&assignment, input.AssignmentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		return
	}

	if assignment.Status != "pending" {
		c.JSON(http.StatusConflict, gin.H{"error": "Assignment already accepted or rejected"})
		return
	}

	assignment.Status = "accepted"
	assignment.AcceptedDriverID = &input.DriverID

	if err := database.DB.Save(&assignment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to accept assignment"})
		return
	}

	database.DB.Model(&models.Assignment{}).Where("vehicle_id = ? AND start_time = ? AND end_time = ? AND id != ?",
		assignment.VehicleID, assignment.StartTime, assignment.EndTime, assignment.ID).
		Updates(map[string]interface{}{"status": "rejected"})

	c.JSON(http.StatusOK, gin.H{"data": assignment})
}

func RejectAssignment(c *gin.Context) {
	var input struct {
		DriverID     uint `json:"driver_id"`
		AssignmentID uint `json:"assignment_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var assignment models.Assignment
	if err := database.DB.First(&assignment, input.AssignmentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		return
	}

	// Log the current status before checking for conflicts
	log.Println("Current assignment status:", assignment.Status)

	if assignment.Status != "pending" {
		c.JSON(http.StatusConflict, gin.H{"error": "Assignment already accepted or rejected"})
		return
	}

	assignment.Status = "rejected"

	if err := database.DB.Save(&assignment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject assignment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": assignment})
}

func GetDriverAssignments(c *gin.Context) {
	driverID := c.Param("driver_id")
	log.Println("Fetching assignments for Driver ID:", driverID)

	var assignments []models.Assignment
	if err := database.DB.Where("driver_id = ?", driverID).Find(&assignments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch assignments"})
		return
	}

	if len(assignments) == 0 {
		c.JSON(http.StatusOK, gin.H{"data": []models.Assignment{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": assignments})
}
