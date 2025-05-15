package controllers

import (
	"Cabinetservice/app/db"
	"Cabinetservice/app/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func CreateCabinet(c *gin.Context) {
	var user models.CreateCabinet
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Ошибка принятия email", err)
		c.JSON(http.StatusConflict, gin.H{"error": "Не правильные данные"})
		return
	}
	if user.Secret != os.Getenv("SECRET") {
		log.Println("Не правильный секрет.Ошибка принятия email")
		c.JSON(http.StatusConflict, gin.H{"error": "Не правильные данные"})
		return
	}
	result := db.DB.Create(&models.Setting{Email: user.Email})
	if result.Error != nil {
		log.Println("Ошибка создания пользователя", result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": "пользователь не создан"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
