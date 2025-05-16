package cleanup

import (
	"Streamserver/internal/models"
	"log"
	"time"

	"gorm.io/gorm"
)

func DeleteStaleRooms(db *gorm.DB, threshold time.Duration) error {
	cutoff := time.Now().Add(-threshold)
	res := db.
		Where("state = ? AND created_at < ?", models.RoomStateWaiting, cutoff).
		Delete(&models.Room{})
	if res.Error != nil {
		return res.Error
	}
	if count := res.RowsAffected; count > 0 {
		log.Printf("cleanup: deleted %d stale rooms (older than %s)\n", count, threshold)
	}
	return nil
}
