package controllers

import (
	"log"
	"motorq_backend/database"
	"motorq_backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateDriver(c *gin.Context) {
	var input models.Driver

	log.Println("Incoming Driver Data: ", c.Request.Body)

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error binding JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Status == "" {
		input.Status = "available"
	}

	if err := database.DB.Create(&input).Error; err != nil {
		log.Println("Error creating driver: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create driver"})
		return
	}

	log.Println("Driver created successfully: ", input)
	c.JSON(http.StatusOK, gin.H{"data": input})
}

func GetDrivers(c *gin.Context) {
	var drivers []models.Driver
	result := database.DB.Session(&gorm.Session{PrepareStmt: false}).Find(&drivers)

	if result.Error != nil {
		log.Println("Error fetching drivers:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve drivers"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"data": []models.Driver{}})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": drivers})
	}
}

func SearchDrivers(c *gin.Context) {
	var drivers []models.Driver
	name := c.Query("name")
	phone := c.Query("phone")

	database.DB.Where("name LIKE ?", "%"+name+"%").Or("phone LIKE ?", "%"+phone+"%").Find(&drivers)

	c.JSON(http.StatusOK, gin.H{"data": drivers})
}

func GetLastRide(c *gin.Context) {
	driverID := c.Param("driver_id")

	var lastAssignment models.Assignment
	err := database.DB.Where("driver_id = ?", driverID).Order("end_time DESC").First(&lastAssignment).Error

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"lastRide": "No rides yet"})
		return
	}

	var vehicle models.Vehicle
	database.DB.First(&vehicle, lastAssignment.VehicleID)

	c.JSON(http.StatusOK, gin.H{"lastRide": vehicle.Make + " " + vehicle.Model})
}
