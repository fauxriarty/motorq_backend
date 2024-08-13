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
	r.GET("/drivers/:driver_id/assignments", controllers.GetDriverAssignments)

	r.POST("/assignments", controllers.AssignVehicleToDriver)
	r.POST("/assignments/unassign", controllers.UnassignVehicleFromDriver)

	r.POST("/assignments/accept", controllers.AcceptAssignment)
	r.POST("/assignments/reject", controllers.RejectAssignment)

	r.GET("/drivers/:driver_id/lastRide", controllers.GetLastRide)

	return r
}
