package dto

import (
	"time"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateProfileRequest struct {
	Name  string `json:"name" binding:"omitempty"`
	Email string `json:"email" binding:"omitempty,email"`
	Bio   string `json:"bio" binding:"omitempty"`
}

type BaseResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type RegisterResponse struct {
	BaseResponse
	Data UserResponse `json:"data"`
}

type LoginResponse struct {
	BaseResponse
	Token string       `json:"token"`
	Data  UserResponse `json:"data"`
}

type UserResponse struct {
	UserId string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

type UpdateProfileResponse struct {
	UserId string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Bio    string `json:"bio"`
}

type EmailRequest struct {
	To      string `json:"to" binding:"required"`
	Subject string `json:"subject" binding:"required"`
	Message string `json:"message" binding:"required"`
}

type ForgotPasswordReq struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordReq struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type ResponseError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type BlacklistedToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Token     string    `gorm:"type:text;not null;unique"`
	CreatedAt time.Time
	ExpiresAt time.Time
}
