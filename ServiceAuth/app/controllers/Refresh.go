package controllers

import (
	"ServiceAuth/app/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func RefreshToken(c *gin.Context) {
	refreshtoken, err := c.Cookie("refresh_token")
	if err != nil || refreshtoken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Нет refresh токена"})
		return
	}
	claims, err := utils.ParseRefreshToken(refreshtoken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Невалидный refresh токен"})
		return
	}
	email := claims.Email

	stored, err := utils.RedisClient.Get(utils.Ctx, "refresh_token").Result()
	if err != nil || stored != refreshtoken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Токен не найден или просрочен"})
		return
	}
	newAccessToken, err := utils.GenerateAccessToken(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать access токен"})
		return
	}

	domain := os.Getenv("DOMAIN")
	c.SetCookie("access_token", newAccessToken, 900, "/", domain, false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Access токен обновлён"})
}
