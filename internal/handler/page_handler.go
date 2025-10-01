package handler

import (
	"net/http"

	"github.com/bayuTri-Code/BE-Recipe/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DashboardHandler struct {
	Service *services.DashboardService
}

func NewDashboardHandler(s *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{Service: s}
}
type ErrorResponse struct {
	Error string `json:"error"`
}


// GetDashboard godoc
// @Summary Get user dashboard data
// @Description Get user dashboard information including total recipes and favorites
// @Tags Dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.DashboardDTO
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/page/dashboard [get]
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	data, err := h.Service.GetDashboardData(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}
