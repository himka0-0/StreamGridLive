package handlers

import (
	"Cabinetservice/app/db"
	"Cabinetservice/app/models"
	"Cabinetservice/app/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strconv"
)

func GoStrimPage(c *gin.Context) {
	email := c.GetString("email")
	if email == "" {
		log.Println("Проблема в контексте админа")
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "Пользователь не авторизован"})
		return
	}
	var settings models.Setting
	if err := db.DB.Where("email=?", email).First(&settings).Error; err != nil {
		log.Println("ошибка поиска настроек трансляции пользователя", err)
		c.JSON(http.StatusConflict, gin.H{"error": "Отредактируйте настройки на вкладке настроить кабинет"})
		return
	}
	invitationLink := utils.GenerateLink()
	invitationLink = os.Getenv("STRIM_URL") + "/id=" + strconv.Itoa(int(settings.ID)) + invitationLink
	c.HTML(http.StatusOK, "StrimBoard.html", gin.H{
		"invitationLink": invitationLink,
		"setting":        settings,
		"password":       settings.Password,
	})
}
