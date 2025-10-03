package utils

import (
	"errors"
	"net/smtp"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretReset = []byte(os.Getenv("ACCESS_TOKEN_RESET"))

func GenerateTokenReset(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"reset":   true,
		"exp": time.Now().Add(30 * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretReset)
}

func ValidateTokenReset(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecretReset, nil
	})

	if err != nil {
		return nil, errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}


func SendEmail(to, subject, body string) error {
	from := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")

	addr := host + ":" + port

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	auth := smtp.PlainAuth("", from, pass, host)

	if err := smtp.SendMail(addr, auth, from, []string{to}, []byte(msg)); err != nil {
		return err
	}
	return nil
}