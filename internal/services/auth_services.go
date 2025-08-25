package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bayuTri-Code/BE-Recipe/database"
	"github.com/bayuTri-Code/BE-Recipe/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Claims struct {
	ID   uuid.UUID `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

var jwtSecret []byte

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("File .env is not found")
	}
	secret := os.Getenv("ACCESS_TOKEN_SECRET")
	if secret == "" {
		log.Fatal("ACCESS_TOKEN_SECRET not set in environment")
	}
	jwtSecret = []byte(secret)
}

func RegisterServices(ctx context.Context, Name, Email, Password string) (*models.User, error) {
	db := database.Db
	Email = strings.ToLower(strings.TrimSpace(Email))

	var existing models.User
	if err := db.WithContext(ctx).Where("email = ?", Email).First(&existing).Error; err == nil {
		return nil, errors.New("user already exists")
	}

	registerUser := &models.User{
		ID:       uuid.New(),
		Name:     Name,
		Email:    Email,
		Password: Password,
	}
	if err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(registerUser).Error; err != nil {
			return fmt.Errorf("failed to create user: %v", err)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction failed: %v", err)
	}
	return registerUser, nil
}


func LoginServices(email, password string) (string, error) {
	var user models.User

	result := database.Db.Where("email = ? AND password = ?", email, password).First(&user)
	if result.Error != nil {
		return "", errors.New("invalid email or password")
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		ID:   user.ID,
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", errors.New("could not create token")
	}

	return tokenString, nil
}


func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := database.Db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}