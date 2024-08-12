package routes

import (
	"motorq_backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/drivers", controllers.CreateDriver)
	r.GET("/drivers", controllers.GetDrivers)
	r.GET("/drivers/search", controllers.SearchDrivers)

	r.POST("/vehicles", controllers.CreateVehicle)
	r.GET("/vehicles", controllers.GetVehicles)

	r.POST("/assignments", controllers.AssignVehicleToDriver)
	r.POST("/assignments/unassign", controllers.UnassignVehicleFromDriver)

	return r
}
