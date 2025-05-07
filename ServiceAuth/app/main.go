package main

import (
	"ServiceAuth/app/db"
	"ServiceAuth/app/models"
	"ServiceAuth/app/router"
	"ServiceAuth/app/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}
	db.InitDB()
	utils.InitRedis()
	db.DB.AutoMigrate(&models.User{})

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	router.SetupRoutes(r)
	r.Run(":8081")
}
