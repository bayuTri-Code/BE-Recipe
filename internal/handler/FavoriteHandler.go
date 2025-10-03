package handler

import (
	"net/http"

	"github.com/bayuTri-Code/BE-Recipe/internal/services"
	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	Service *services.FavoriteService
}

func NewFavoriteHandler(s *services.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{Service: s}
}

// Favorite godoc
// @Summary      post recipe favorite
// @Description  This endpoint will add the recipe to user's favorites if it is not already favorited. 
// @Description  If the recipe is already in favorites, it will remove it instead.
// @Tags         Favorites
// @Produce      json
// @Security     BearerAuth
// @Param        recipe_id   path      string  true  "Recipe ID"
// @Success      200  {object}  map[string]string "message: 'Added to favorites' or 'Removed from favorites'"
// @Failure      400  {object}  map[string]string "Invalid recipe ID"
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /api/recipes/{recipe_id}/favorites [post]
func (h *FavoriteHandler) AddFavoriteHandler(c *gin.Context) {
	recipeID := c.Param("recipe_id")
	if recipeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "recipe_id is required"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	added, err := h.Service.AddFavoriteService(userID.(string), recipeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if added {
		c.JSON(http.StatusOK, gin.H{"message": "recipe added to favorites"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "recipe removed from favorites"})
	}
	 c.JSON(http.StatusOK, gin.H{
        "recipe_id": recipeID,
        "favorited": added, 
    })
}


// // RemoveFavorite godoc
// // @Summary      Remove recipe from favorites
// // @Description  Explicitly remove recipe from user's favorites
// // @Tags         Favorites
// // @Produce      json
// // @Security     BearerAuth
// // @Param        recipe_id   path      string  true  "Recipe ID"
// // @Success      204  {object}  nil
// // @Failure      400  {object}  map[string]string
// // @Failure      401  {object}  map[string]string
// // @Failure      500  {object}  map[string]string
// // @Router       /api/recipes/{recipe_id}/favorites [delete]
// func (h *FavoriteHandler) RemoveFavorite(c *gin.Context) {
// 	recipeID := c.Param("recipe_id")
// 	if recipeID == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "recipe_id is required"})
// 		return
// 	}

// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 		return
// 	}

// 	if err := h.Service.RemoveFavorite(userID.(string), recipeID); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove favorite"})
// 		return
// 	}

// 	c.Status(http.StatusNoContent)
// }



// GetAllFavorites godoc
// @Summary      Get all favorite recipes by user
// @Description  Get all favorite recipes of the logged-in user
// @Tags         Favorites
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.Favorite
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/recipes/favorites [get]
func (h *FavoriteHandler) GetAllFavorites(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	favorites, err := h.Service.GetAllFavorites(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get favorites: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, favorites)
}
