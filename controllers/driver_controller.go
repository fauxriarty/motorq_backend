package controllers

import (
	"motorq_backend/database"
	"motorq_backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateDriver(c *gin.Context) {
	var input models.Driver
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	driver := models.Driver{
		Name:      input.Name,
		Email:     input.Email,
		Phone:     input.Phone,
		Location:  input.Location,
		WorkHours: input.WorkHours,
	}

	database.DB.Create(&driver)

	c.JSON(http.StatusOK, gin.H{"data": driver})
}

func GetDrivers(c *gin.Context) {
	var drivers []models.Driver
	database.DB.Find(&drivers)

	c.JSON(http.StatusOK, gin.H{"data": drivers})
}

func SearchDrivers(c *gin.Context) {
	var drivers []models.Driver
	name := c.Query("name")
	phone := c.Query("phone")

	database.DB.Where("name LIKE ?", "%"+name+"%").Or("phone LIKE ?", "%"+phone+"%").Find(&drivers)

	c.JSON(http.StatusOK, gin.H{"data": drivers})
}
