package main

import (
	"Streamserver/internal/api"
	"Streamserver/internal/cleanup"
	"Streamserver/internal/db"
	"Streamserver/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"time"
)

type createRoomRequest struct {
	InvitationLink string `json:"invitationLink" binding:"required,url"`
	Tool           string `json:"tool" binding:"required"`
	Permissions    string `json:"permissions" binding:"required"`
	Password       string `json:"password"` // необязательно
}

type createRoomResponse struct {
	ID string `json:"id"`
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: добавить проверку origin, если нужно
		return true
	},
}

func main() {

	gormDB, err := db.Init()
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			if err := cleanup.DeleteStaleRooms(gormDB, 2*time.Minute); err != nil {
				log.Printf("cleanup error: %v\n", err)
			}
		}
	}()

	r := gin.Default()

	r.GET("/ws", api.WSHandler(gormDB))
	r.POST("/api/rooms", func(c *gin.Context) {
		var req createRoomRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		roomID := uuid.New()
		room := models.Room{
			ID:             roomID,
			InvitationLink: req.InvitationLink,
			Tool:           req.Tool,
			Permissions:    req.Permissions,
			Password:       req.Password,
			State:          "waiting",
		}

		if err := gormDB.Create(&room).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create room"})
			return
		}

		c.JSON(http.StatusOK, createRoomResponse{ID: roomID.String()})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("listening on :%s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
