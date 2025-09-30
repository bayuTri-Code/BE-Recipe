package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bayuTri-Code/BE-Recipe/internal/models"
	"github.com/bayuTri-Code/BE-Recipe/internal/services"
	"github.com/bayuTri-Code/BE-Recipe/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// @Summary Register user
// @Description Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Register request"
// @Success 201 {object} models.RegisterResponse
// @Failure 400 {object} models.ResponseError
// @Router /auth/register [post]
func RegisterHandler(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid Input: "+err.Error())
		return
	}

	user, err := services.RegisterServices(c, req.Name, req.Email, req.Password)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusCreated, models.RegisterResponse{
		BaseResponse: models.BaseResponse{
			Success: true,
			Message: "User created successfully",
		},
		Data: models.UserResponse{
			UserId:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
		},
	})
}

// @Summary Login user
// @Description Authenticate user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login request"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.ResponseError
// @Router /auth/login [post]
func LoginHandler(c *gin.Context) {
	var req models.LoginRequest
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

	c.JSON(http.StatusOK, models.LoginResponse{
		BaseResponse: models.BaseResponse{
			Success: true,
			Message: "Login success",
		},
		Token: token,
		Data: models.UserResponse{
			UserId: user.ID.String(),
			Name:   user.Name,
			Email:  user.Email,
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
