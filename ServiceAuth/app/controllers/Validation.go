package controllers

import (
	"ServiceAuth/app/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func ValidationToken(c *gin.Context) {
	autHeader := c.GetHeader("Authorization")
	if autHeader == "" || !strings.HasPrefix(autHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"valid": false})
		return
	}

	tokenStr := strings.TrimPrefix(autHeader, "Bearer ")

	claims, err := utils.ParseRefreshToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"valid": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid": true,
		"email": claims.Email,
	})
}
