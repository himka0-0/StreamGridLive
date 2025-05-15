package main

import (
	"Cabinetservice/app/db"
	"Cabinetservice/app/models"
	"Cabinetservice/app/routers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}
	db.InitDB()
	db.DB.AutoMigrate(&models.Setting{})
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	routers.SetupRoutes(r)
	certFile := "certs/_wildcard.dev.local.pem"
	keyFile := "certs/_wildcard.dev.local-key.pem"
	log.Fatal(r.RunTLS(":8082", certFile, keyFile))
}
