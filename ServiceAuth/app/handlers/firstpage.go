package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Firstpage(c *gin.Context) {
	c.HTML(http.StatusOK, "firstPage.html", gin.H{})
}
