package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"os"
	"time"
)

var jwtSecret = []byte(os.Getenv("ACCESS_TOKEN_SECRET"))

func GenerateJWT(ID uuid.UUID, Name, Email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": ID.String(),
		"name":    Name,
		"email":   Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func VerifyJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
}