package utils

import (
	"encoding/hex"
	"math/rand"
)

func GenerateLink() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
