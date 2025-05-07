package utils

import (
	"ServiceAuth/app/models"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func GenerateAccessToken(email string) (string, error) {
	jwtkey := []byte(os.Getenv("JWT_KEY"))

	claims := models.Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtkey)
}
