package controllers

import (
	"ServiceAuth/app/db"
	"ServiceAuth/app/models"
	"ServiceAuth/app/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func Registration(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Ошибка парсингда данных при регистрации пользователя", err)
		c.JSON(http.StatusConflict, gin.H{"error": "Не правильные даннные"})
		return
	}
	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("ошибка хеширования пароля", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка сервера"})
		return
	}
	verifyToken := utils.GenerationToken()
	user.Verification_token = verifyToken
	user.Password = string(hashPass)

	err = db.DB.Create(&user).Error
	if err != nil {
		log.Println("Ошибка сохранения ", err)
		c.JSON(http.StatusConflict, gin.H{"error": "Данные не сохранены"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Регистрация успешна!Нобходимо подтвердить почту",
	})
	
	err = utils.PublishVerificationEmail(user.Email, verifyToken)
	if err != nil {
		log.Println("ошибка публикации в rabbit")
	}
}
