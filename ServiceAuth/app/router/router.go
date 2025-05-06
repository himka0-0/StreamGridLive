package router

import (
	"ServiceAuth/app/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/", handlers.Firstpage)

	r.GET("/login", handlers.Login)

	r.GET("/registration", handlers.Registration)

	r.GET("/recovery", handlers.Recovery)

	r.GET("/messagemail", handlers.Messagemail)

	r.GET("/recoverypassword")
}
