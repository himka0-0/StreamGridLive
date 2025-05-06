package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Recovery(c *gin.Context) {
	c.HTML(http.StatusOK, "RecoveryPass.html", gin.H{})
}

func Messagemail(c *gin.Context) {
	c.HTML(http.StatusOK, "Stopmail.html", gin.H{})
}
