package controllers

import (
	"ServiceAuth/app/db"
	"ServiceAuth/app/models"
	"ServiceAuth/app/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func Login(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Ошибка парсинга данных,Login", err)
		c.JSON(http.StatusConflict, gin.H{"error": "Введены не правильные данные"})
		return
	}
	var user models.User
	err := db.DB.Where("email=?", input.Email).First(&user).Error
	if err != nil {
		log.Println("Авторизация по не существуещим данным,Login", err)
		c.JSON(http.StatusConflict, gin.H{"error": "Не правильный логин или пароль"})
		return
	}
	if !user.Verify_mail {
		c.JSON(http.StatusConflict, gin.H{"error": "Для входа подтвердите почту"})
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Не правильный логин или пароль"})
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	var accessToken, refreshToken string
	var accessErr, refreshErr error

	go func() {
		defer wg.Done()
		accessToken, accessErr = utils.GenerateAccessToken(user.Email)
	}()
	go func() {
		defer wg.Done()
		refreshToken, refreshErr = utils.GenerateRefreshToken(user.Email)
	}()

	wg.Wait()

	if accessErr != nil || refreshErr != nil {
		log.Println("Ошибка генерации токенов, Login:", accessErr, refreshErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка авторизации"})
		return
	}
	utils.RedisClient.Set(utils.Ctx, "refresh_"+user.Email, refreshToken, 7*24*time.Hour)

	domain := os.Getenv("DOMAIN")

	c.SetCookie("access_token", accessToken, 900, "/", domain, false, true)
	c.SetCookie("refresh_token", refreshToken, int(7*24*time.Hour.Seconds()), "/", domain, false, true)

	location := fmt.Sprintf("http://%s/mycabinet", os.Getenv("CABINET_SERVICE"))
	c.Redirect(http.StatusFound, location)
}
