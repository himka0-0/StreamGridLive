package controllers

import (
	"Cabinetservice/app/db"
	"Cabinetservice/app/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Setting(c *gin.Context) {
	email := c.GetString("email")
	if email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}
	var input models.Setting
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Не возможно запарсить данные", err)
		c.JSON(http.StatusConflict, gin.H{"error": "Не правильные данные"})
		return
	}
	result := db.DB.Model(&models.Setting{}).Where("email = ?", email).Updates(models.Setting{
		Tool:     input.Tool,
		Rule:     input.Rule,
		Password: input.Password,
	})

	if result.Error != nil {
		log.Println("Ошибка при сохранении данных настроек:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера при сохранении"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Настройки не найдены для обновления"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Данные сохранены"})
}
