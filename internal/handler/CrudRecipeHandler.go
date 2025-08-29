package handler

import (
	"net/http"

	dto "github.com/bayuTri-Code/BE-Recipe/internal/DTO"
	"github.com/bayuTri-Code/BE-Recipe/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RecipeHandler struct {
	Service *services.RecipeService
}

func NewRecipeHandler(s *services.RecipeService) *RecipeHandler {
	return &RecipeHandler{Service: s}
}

// CreateRecipe godoc
// @Summary Create a new recipe by the authenticated user
// @Description Create a new recipe with title, description, category, prep_time, cook_time and ingredients and steps for the logged-in user
// @Tags Recipes
// @Accept json
// @Produce json
// @Param recipe body dto.CreateRecipeRequest true "Recipe Data"
// @Success 201 {object} models.Recipe
// @Failure 400 {object} map[string]string
// @Router /api/recipes [post]
func (h *RecipeHandler) CreateRecipe(c *gin.Context) {
	var req dto.CreateRecipeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	
	userID := c.GetString("userID")
	uid, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
		return
	}

	
	res, err := h.Service.CreateRecipe(req, uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}



// GetRecipes godoc
// @Summary Get all recipes for all users
// @Description Retrieve all available recipes
// @Tags Recipes
// @Produce json
// @Success 200 {array} models.Recipe
// @Router /api/recipes [get]
func (h *RecipeHandler) GetAllRecipes(c *gin.Context) {
	list, err := h.Service.GetAllRecipes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch recipes"})
		return
	}
	c.JSON(http.StatusOK, list)
}


// GetRecipeByID godoc
// @Summary Get recipe by ID
// @Description Retrieve a recipe by its ID
// @Tags Recipes
// @Produce json
// @Param id path string true "Recipe ID"
// @Success 200 {object} models.Recipe
// @Failure 404 {object} map[string]string
// @Router /api/recipes/{id} [get]
func (h *RecipeHandler) GetRecipeByID(c *gin.Context) {
	id := c.Param("id")
	res, err := h.Service.GetRecipeByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
		return
	}
	c.JSON(http.StatusOK, res)
}


// GetMyRecipes godoc
// @Summary Get my recipes 
// @Description Get all recipes created by the authenticated user
// @Tags Recipes
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Recipe
// @Failure 401 {object} map[string]string
// @Router /api/myrecipes [get]
func (h *RecipeHandler) GetMyRecipes(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	recipes, err := h.Service.GetRecipesByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recipes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"recipes": recipes})
}


// UpdateRecipe godoc
// @Summary Update recipe
// @Description Update recipe details by ID
// @Tags Recipes
// @Accept json
// @Produce json
// @Param id path string true "Recipe ID"
// @Param recipe body models.Recipe true "Updated Recipe"
// @Success 200 {object} models.Recipe
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/recipes/{id} [put]
func (h *RecipeHandler) UpdateRecipe(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateRecipeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.Service.UpdateRecipe(id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// DeleteRecipe godoc
// @Summary Delete recipe
// @Description Delete a recipe by ID
// @Tags Recipes
// @Produce json
// @Param id path string true "Recipe ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/recipes/{id} [delete]
func (h *RecipeHandler) DeleteRecipe(c *gin.Context) {
	id := c.Param("id")
	if err := services.DeleteRecipeService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete recipe"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "recipe deleted"})
}


