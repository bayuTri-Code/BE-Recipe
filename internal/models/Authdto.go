package models

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
    UserId    string `json:"user_id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type ResponseError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
