package handlers

import (
	"ServiceAuth/app/db"
	"ServiceAuth/app/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func RegistrationPage(c *gin.Context) {
	c.HTML(http.StatusOK, "Registration.html", gin.H{})
}

func Verificationmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Токен отсутствует"})
		return
	}
	var user models.User
	if err := db.DB.Where("verification_token=?", token).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Недействительный токен"})
		return
	}

	user.Verify_mail = true
	user.Verification_token = ""
	err := db.DB.Save(&user).Error
	if err != nil {
		log.Println("Ошибка сохранения подтверждения почты", err)
		c.JSON(http.StatusConflict, gin.H{"error": "Ошибка сохранения"})
		return
	}
	c.HTML(http.StatusOK, "Verificationmail.html", gin.H{})
}
