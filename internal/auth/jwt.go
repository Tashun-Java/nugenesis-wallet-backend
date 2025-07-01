package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var jwtSecret = []byte("supersecretkey") // Ideally from an .env

func GenerateJWT(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
