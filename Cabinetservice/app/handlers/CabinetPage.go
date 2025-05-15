package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func CabinetPage(c *gin.Context) {
	username := c.GetString("username")
	if username == "" {
		log.Println("Проблема в контексте админа")
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "Пользователь не авторизован"})
		return
	}
	c.HTML(http.StatusOK, "Cabinet.html", gin.H{"username": username})
}
