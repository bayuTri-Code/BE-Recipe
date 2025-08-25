package handler

import (
	"net/http"

	dto "github.com/bayuTri-Code/BE-Recipe/internal/DTO"
	"github.com/bayuTri-Code/BE-Recipe/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FavoriteHandler struct {
	Service *services.FavoriteService
}

func NewFavoriteHandler(s *services.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{Service: s}
}

// AddFavorite godoc
// @Summary Add recipe to favorites
// @Description Add a recipe to user's favorite list
// @Tags Favorites
// @Accept json
// @Produce json
// @Param favorite body models.Favorite true "Favorite Data"
// @Success 201 {object} models.Favorite
// @Failure 400 {object} map[string]string
// @Router /api/{user_id}/favorites [post]
func (h *FavoriteHandler) AddFavorite(c *gin.Context) {
	recipeID := c.Param("id")
	var req dto.AddFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.UserID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}
	if err := h.Service.AddFavorite(recipeID, req.UserID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}


// DeleteFavorite godoc
// @Summary Remove recipe from favorites
// @Description Delete a favorite by ID
// @Tags Favorites
// @Produce json
// @Param id path string true "Favorite ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/{id}/favorites/{user_id} [delete]
func (h *FavoriteHandler) RemoveFavorite(c *gin.Context) {
	recipeID := c.Param("id")
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}
	if err := h.Service.RemoveFavorite(recipeID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove favorite"})
		return
	}
	c.Status(http.StatusNoContent)
}



// GetAllFavorites godoc
// @Summary      Get all favorite recipes by user
// @Description  Get all favorite recipes of the logged in user
// @Tags         Favorites
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.Favorite
// @Failure      401  {object}  map[string]string{error=string}
// @Failure      500  {object}  map[string]string{error=string}
// @Router       /api/favorites [get]
func GetAllFavorites(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	favorites, err := services.GetAllFavoritesService(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get favorites"})
		return
	}

	c.JSON(http.StatusOK, favorites)
}