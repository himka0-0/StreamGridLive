package cleanup_test

import (
	"Streamserver/internal/cleanup"
	"Streamserver/internal/models"
	"github.com/google/uuid"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestDeleteStaleRooms(t *testing.T) {
	// in-memory БД
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.Room{}))

	// создаём "свежую" и "старую" комнаты
	fresh := models.Room{
		ID:        uuid.New(),
		State:     models.RoomStateWaiting,
		CreatedAt: time.Now(),
	}
	stale := models.Room{
		ID:        uuid.New(),
		State:     models.RoomStateWaiting,
		CreatedAt: time.Now().Add(-5 * time.Minute),
	}
	db.Create(&fresh)
	db.Create(&stale)

	// выполняем очистку с порогом 2 минуты
	err = cleanup.DeleteStaleRooms(db, 2*time.Minute)
	require.NoError(t, err)

	var count int64
	db.Model(&models.Room{}).Count(&count)
	require.Equal(t, int64(1), count, "должна остаться только свежая комната")
}
