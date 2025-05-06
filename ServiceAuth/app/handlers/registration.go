package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Registration(c *gin.Context) {
	c.HTML(http.StatusOK, "Registration.html", gin.H{})
}
