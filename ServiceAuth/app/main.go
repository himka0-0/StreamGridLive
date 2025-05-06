package main

import (
	"ServiceAuth/app/db"
	"ServiceAuth/app/router"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	router.SetupRoutes(r)
	r.Run(":8081")
}
