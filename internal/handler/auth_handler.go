package handler

import (
	"net/http"

	"github.com/bayuTri-Code/Auth-Services/internal/models"
	"github.com/bayuTri-Code/Auth-Services/internal/services"
	"github.com/bayuTri-Code/Auth-Services/internal/utils"
	"github.com/gin-gonic/gin"
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
			ID:    user.ID.String(),
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
			ID:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
		},
	})
}
