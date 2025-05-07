package router

import (
	"ServiceAuth/app/controllers"
	"ServiceAuth/app/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/", handlers.Firstpage)

	r.GET("/login", handlers.LoginPage)
	r.POST("/login", controllers.Login)

	r.GET("/registration", handlers.RegistrationPage)
	r.POST("/registration", controllers.Registration)

	r.GET("/Verificationmail", handlers.Verificationmail)

	r.GET("/recovery", handlers.Recovery)

	r.GET("/messagemail", handlers.Messagemail)

	r.GET("/recoverypassword")

	r.POST("/validation", controllers.ValidationToken)
}
