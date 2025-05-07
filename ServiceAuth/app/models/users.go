package models

import "github.com/golang-jwt/jwt/v5"

type User struct {
	ID                 uint   `gorm:"primary_key"`
	Name               string `gorm:"size 100;not null;unique" json:"name"`
	Email              string `gorm:"size 100;not null;unique" json:"email"`
	Password           string `gorm:"not null" json:"password"`
	Verification_token string `gorm:"size 100"`
	Verify_mail        bool   `gorm:"not null;default:false"`
}
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}
