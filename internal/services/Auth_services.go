package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bayuTri-Code/BE-Recipe/database"
	dto "github.com/bayuTri-Code/BE-Recipe/internal/DTO"
	"github.com/bayuTri-Code/BE-Recipe/internal/models"
	"github.com/bayuTri-Code/BE-Recipe/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Claims struct {
	UserId uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Bio    string    `json:"bio"`
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

func saveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.ReadFrom(src)
	return err
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
		Bio:    user.Bio,
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

func UpdateProfile(userID uuid.UUID, name, email, bio string, avatarFile, bannerFile *multipart.FileHeader) (*dto.UpdateProfileResponse, error) {
	var user models.User
	if err := database.Db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	if name != "" {
		user.Name = name
	}
	if email != "" {
		user.Email = email
	}
	if bio != "" {
		user.Bio = bio
	}

	apiImagePath := os.Getenv("API_IMAGE_PATH")


	if avatarFile != nil {
		uploadDir := "public/profile_storage"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return nil, errors.New("failed to create avatar storage directory")
		}

		if user.Avatar != "" {
			_ = os.Remove(filepath.Join(uploadDir, filepath.Base(user.Avatar)))
		}
		
		if avatarFile.Size > 2*1024*1024 {
			return nil, errors.New("avatar file size must not exceed 2MB")
		}

		fileExt := filepath.Ext(avatarFile.Filename)
		newFileName := fmt.Sprintf("%s%s", uuid.NewString(), fileExt)
		filePath := filepath.Join(uploadDir, newFileName)

		if err := saveUploadedFile(avatarFile, filePath); err != nil {
			return nil, errors.New("failed to save avatar file")
		}

		user.Avatar = fmt.Sprintf("%s/profile-storage/%s", apiImagePath, newFileName)
	}

	if bannerFile != nil {
		uploadDir := "public/profile_banner"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return nil, errors.New("failed to create banner storage directory")
		}

		if user.Banner != "" {
			_ = os.Remove(filepath.Join(uploadDir, filepath.Base(user.Banner)))
		}

		if bannerFile.Size > 2*1024*1024 {
			return nil, errors.New("avatar file size must not exceed 2MB")
		}

		fileExt := filepath.Ext(bannerFile.Filename)
		newFileName := fmt.Sprintf("%s%s", uuid.NewString(), fileExt)
		filePath := filepath.Join(uploadDir, newFileName)

		if err := saveUploadedFile(bannerFile, filePath); err != nil {
			return nil, errors.New("failed to save banner file")
		}

		user.Banner = fmt.Sprintf("%s/profile-banner/%s", apiImagePath, newFileName)
	}

	if err := database.Db.Save(&user).Error; err != nil {
		return nil, errors.New("failed to update user")
	}

	return &dto.UpdateProfileResponse{
		UserId: user.ID.String(),
		Name:   user.Name,
		Email:  user.Email,
		Bio:    user.Bio,
		Avatar: user.Avatar,
		Banner: user.Banner,
	}, nil
}

func BlacklistToken(tokenString string, expiresAt time.Time) error {
	fmt.Println("Blacklisting token:", tokenString)

	blacklisted := dto.BlacklistedToken{
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
	fmt.Println("cek token di blacklist:", tokenString)

	var token dto.BlacklistedToken
	err := database.Db.Where("token = ?", strings.TrimSpace(tokenString)).First(&token).Error

	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	fmt.Println("Token terblacklist:", token.Token)
	return true, nil
}

func ForgotPassword(email string) (string, error) {
	var user models.User
	if err := database.Db.Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("email not found")
	}

	token, err := utils.GenerateTokenReset(user.ID.String())
	if err != nil {
		return "", err
	}

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", os.Getenv("APP_URL"), token)

	if os.Getenv("APP_ENV") == "development" {
		return token, nil
	}

	subject := "Reset Your Password"
	body := fmt.Sprintf("Klik link berikut untuk reset password:\n\n%s", resetLink)
	if err := utils.SendEmail(user.Email, subject, body); err != nil {
		return "", err
	}

	return "Reset link sent to your email", nil
}

func ResetPassword(token, newPassword string) (bool, error) {
	claims, err := utils.ValidateTokenReset(token)
	if err != nil {
		return false, errors.New("invalid or expired token")
	}

	userID := claims["user_id"].(string)

	if err := database.Db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("password", newPassword).Error; err != nil {
		return false, err
	}

	return true, nil
}
