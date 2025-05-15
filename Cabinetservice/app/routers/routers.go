package routers

import (
	"Cabinetservice/app/controllers"
	"Cabinetservice/app/handlers"
	"Cabinetservice/app/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	userRoute := r.Group("/mycabinet")
	userRoute.Use(middlewares.AuthMiddleware())
	{
		userRoute.GET("/", handlers.CabinetPage)

		userRoute.GET("/settings", handlers.SettingPage)
		userRoute.POST("/settings", controllers.Setting)

		userRoute.GET("/gostrim", handlers.GoStrimPage)
	}
	r.POST("/createcabinet", controllers.CreateCabinet)
}
