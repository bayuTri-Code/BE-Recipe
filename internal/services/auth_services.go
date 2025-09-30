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
	UserId uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
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
		UserId: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()), 
			ID:        uuid.NewString(),
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
func BlacklistToken(tokenString string, expiresAt time.Time) error {
    fmt.Println("Blacklisting token:", tokenString)

    blacklisted := models.BlacklistedToken{
        ID:        uuid.New(),
        Token:     strings.TrimSpace(tokenString),
        CreatedAt: time.Now(),
        ExpiresAt: expiresAt,
    }

    if err := database.Db.Create(&blacklisted).Error; err != nil {
        fmt.Printf("Gagal blacklist: %v\n", err)
        return fmt.Errorf("failed to blacklist token: %v", err)
    }

    fmt.Println("Token berhasil di-blacklist:", tokenString)
    return nil
}

func IsTokenBlacklisted(tokenString string) (bool, error) {
    fmt.Println("Mengecek token di blacklist:", tokenString)

    var token models.BlacklistedToken
    err := database.Db.Where("token = ?", strings.TrimSpace(tokenString)).First(&token).Error

    if err == gorm.ErrRecordNotFound {
        return false, nil
    }
    if err != nil {
        return false, err
    }

    fmt.Println("Token ditemukan di blacklist:", token.Token)
    return true, nil
}
