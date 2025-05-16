package api

import (
	"Streamserver/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRouter() *gin.Engine {
	// in-memory SQLite
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&models.Room{})

	router := gin.Default()
	router.POST("/api/rooms", func(c *gin.Context) {
		var req struct {
			InvitationLink string `json:"invitationLink" binding:"required,url"`
			Tool           string `json:"tool" binding:"required"`
			Permissions    string `json:"permissions" binding:"required"`
			Password       string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		id := uuid.New().String()
		room := models.Room{
			ID:             uuid.MustParse(id),
			InvitationLink: req.InvitationLink,
			Tool:           req.Tool,
			Permissions:    req.Permissions,
			Password:       req.Password,
			State:          "waiting",
		}
		db.Create(&room)
		c.JSON(http.StatusOK, gin.H{"id": id})
	})
	return router
}

func TestCreateRoom(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	body := `{
        "invitationLink":"https://video.example.com/room/abc123",
        "tool":"screen-share",
        "permissions":"moderator",
        "password":"secret123"
    }`
	req := httptest.NewRequest(http.MethodPost, "/api/rooms", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp["id"], "должен вернуть непустой id")
}
