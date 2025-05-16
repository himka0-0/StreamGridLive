package models

import (
	"github.com/google/uuid"
	"time"
)

type Room struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	InvitationLink string    `gorm:"size:255;not null"`
	Tool           string    `gorm:"size:100;not null"`
	Permissions    string    `gorm:"size:100;not null"`
	Password       string    `gorm:"size:100"` // при необходимости можно убрать nullable
	State          string    `gorm:"size:20;not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

const (
	RoomStateWaiting = "waiting"
	RoomStateActive  = "active"
	RoomStateClosed  = "closed"
)
