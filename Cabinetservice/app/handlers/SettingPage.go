package handlers

import (
	"Cabinetservice/app/db"
	"Cabinetservice/app/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func SettingPage(c *gin.Context) {
	email := c.GetString("email")
	if email == "" {
		log.Println("Проблема в контексте админа")
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "Пользователь не авторизован"})
		return
	}
	var setting models.Setting
	if err := db.DB.Where("email=?", email).First(&setting).Error; err != nil {
		log.Println("ошибка поиска настроек пользователя", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка вашего email,обратитесь в поддержку"})
		return
	}
	c.HTML(http.StatusOK, "Settings.html", gin.H{
		"setting": setting,
	})
}
