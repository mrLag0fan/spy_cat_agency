package routes

import (
	"main/internal/handlers"
	"main/internal/repositories"

	"github.com/gin-gonic/gin"
)

func SetupRouter(catRepo *repositories.CatRepository, missionRepo *repositories.MissionRepository) *gin.Engine {
	r := gin.Default()

	catHandler := handlers.NewCatHandler(catRepo)
	missionHandler := handlers.NewMissionHandler(missionRepo)

	catRoutes := r.Group("/cat")
	{
		catRoutes.POST("", catHandler.CreateCat)
		catRoutes.GET("", catHandler.GetAllCats)
		catRoutes.GET("/:id", catHandler.GetCatByID)
		catRoutes.PUT("/:id/salary", catHandler.UpdateCatSalary)
		catRoutes.DELETE("/:id", catHandler.DeleteCat)
	}

	missionRoutes := r.Group("/mission")
	{
		missionRoutes.POST("", missionHandler.CreateMission)
		missionRoutes.DELETE("/:id", missionHandler.DeleteMission)
		missionRoutes.PUT("/:id/complete", missionHandler.CompleteMission)
		missionRoutes.PUT("/targets/:target_id/notes", missionHandler.UpdateTargetNotes)
		missionRoutes.PUT("/targets/:target_id/complete", missionHandler.MarkTargetAsComplete)
		missionRoutes.DELETE("/targets/:target_id", missionHandler.DeleteTarget)
		missionRoutes.POST("/:id/targets", missionHandler.AddTarget)
		missionRoutes.POST("/:id/assign-cat", missionHandler.AssignCatToMission)
		missionRoutes.GET("", missionHandler.GetAllMissions)
		missionRoutes.GET("/:id", missionHandler.GetMissionByID)
	}

	return r
}
