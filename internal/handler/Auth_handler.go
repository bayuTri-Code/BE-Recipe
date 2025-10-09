package handler

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	dto "github.com/bayuTri-Code/BE-Recipe/internal/DTO"
	"github.com/bayuTri-Code/BE-Recipe/internal/services"
	"github.com/bayuTri-Code/BE-Recipe/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// @Summary Register user
// @Description Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register request"
// @Success 201 {object} dto.RegisterResponse
// @Failure 400 {object} dto.ResponseError
// @Router /auth/register [post]
func RegisterHandler(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid Input: "+err.Error())
		return
	}

	user, err := services.RegisterServices(c, req.Name, req.Email, req.Password)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusCreated, dto.RegisterResponse{
		BaseResponse: dto.BaseResponse{
			Success: true,
			Message: "User created successfully",
		},
		Data: dto.UserResponse{
			UserId: user.ID.String(),
			Name:   user.Name,
			Email:  user.Email,
		},
	})
}

// @Summary Login user
// @Description Authenticate user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ResponseError
// @Router /auth/login [post]
func LoginHandler(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid Input: "+err.Error())
		return
	}

	user, err := services.GetUserByEmail(req.Email)
	if err != nil {
		utils.ResponseError(c, http.StatusUnauthorized, "Email not found")
		return
	}

	token, err := services.LoginServices(req.Email, req.Password)
	if err != nil {
		utils.ResponseError(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{
		BaseResponse: dto.BaseResponse{
			Success: true,
			Message: "Login success",
		},
		Token: token,
		Data: dto.UserResponse{
			UserId: user.ID.String(),
			Name:   user.Name,
			Email:  user.Email,
			Bio:    user.Bio,
		},
	})
}

// LogoutHandler godoc
// @Summary Logout user
// @Description Blacklist token from Authorization header
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string{message=string}
// @Failure 400 {object} map[string]string{error=string}
// @Failure 401 {object} map[string]string{error=string}
// @Failure 500 {object} map[string]string{error=string}
// @Router /logout [post]
func LogoutHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		return
	}
	tokenString := parts[1]

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot parse claims"})
		return
	}

	exp, ok := claims["exp"].(float64)
	var expiresAt time.Time
	if ok {
		expiresAt = time.Unix(int64(exp), 0)
	} else {
		expiresAt = time.Now().Add(24 * time.Hour)
	}

	err = services.BlacklistToken(tokenString, expiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to blacklist token: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful, token blacklisted"})
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the logged-in user's profile, including name, email, bio, avatar image adn banner image.
// @Tags Auth
// @Accept multipart/form-data
// @Produce json
// @Param name formData string false "Name of the user"
// @Param email formData string false "Email address of the user"
// @Param bio formData string false "Short biography"
// @Param avatar formData file false "Avatar image file"
// @Param banner formData file false "banner image file"
// @Success 200 {object} dto.UpdateProfileResponse "Successfully updated profile"
// @Failure 400 {object} map[string]string "Bad request / validation error"
// @Failure 401 {object} map[string]string "Unauthorized / invalid token"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /auth/profile [put]
func UpdateProfileHandler(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	name := c.PostForm("name")
	email := c.PostForm("email")
	bio := c.PostForm("bio")

	avatarFile, _ := c.FormFile("avatar")
	bannerFile, _ := c.FormFile("banner") 

	res, err := services.UpdateProfile(userID, name, email, bio, avatarFile, bannerFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "profile updated successfully",
		"data":    res,
	})
}


// ForgotPasswordHandler godoc
// @Summary Send password reset link
// @Description Generates a password reset token and sends it to the user's email.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ForgotPasswordReq true "Email request"
// @Success 200 {object} map[string]string "Reset token (dev) or message"
// @Failure 400 {object} map[string]string "Invalid email or other error"
// @Router /auth/forgot-password [post]
func ForgotPasswordHandler(c *gin.Context) {
	var req dto.ForgotPasswordReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
		return
	}

	token, err := services.ForgotPassword(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if os.Getenv("APP_ENV") == "development" {
		c.JSON(http.StatusOK, gin.H{"reset_token": token})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Check your email for reset link"})
}

// ResetPasswordHandler godoc
// @Summary Reset user password
// @Description Resets the user's password using a valid reset token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordReq true "Reset Password request"
// @Success 200 {object} map[string]string "Password reset successful"
// @Failure 400 {object} map[string]string "Invalid request or token"
// @Router /auth/reset-password [post]
func ResetPasswordHandler(c *gin.Context) {
	var req dto.ResetPasswordReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	ok, err := services.ResetPassword(req.Token, req.NewPassword)
	if err != nil || !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}
