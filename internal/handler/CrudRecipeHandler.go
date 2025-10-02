package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

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
// @Description Create a new recipe with title, description, category, prep_time, cook_time, ingredients, steps, and thumbnail
// @Tags Recipes
// @Accept multipart/form-data
// @Produce json
// @Param thumbnail formData file false "Thumbnail Image"
// @Success 201 {object} dto.RecipeResponse
// @Failure 400 {object} map[string]string
// @Router /api/recipes [post]
// CreateRecipe godoc
// @Summary Create a new recipe by the authenticated user
// @Description Create a new recipe with title, description, category, prep_time, cook_time, ingredients, steps, and thumbnail
// @Tags Recipes
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Recipe Title"
// @Param description formData string false "Recipe Description"
// @Param category formData string false "Recipe Category"
// @Param prep_time formData int false "Preparation Time"
// @Param cook_time formData int false "Cooking Time"
// @Param servings formData int false "Number of Servings"
// @Param ingredients formData string false "Ingredients JSON Array"
// @Param steps formData string false "Steps JSON Array"
// @Success 201 {object} dto.RecipeResponse
// @Failure 400 {object} map[string]string
// @Router /api/recipes [post]
func (h *RecipeHandler) CreateRecipe(c *gin.Context) {
	userID := c.GetString("userID")
	uid, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
		return
	}

	var req dto.CreateRecipeRequest
	req.Title = c.PostForm("title")
	req.Description = c.PostForm("description")
	req.Category = c.PostForm("category")

	if prepTime := c.PostForm("prep_time"); prepTime != "" {
		if val, err := strconv.Atoi(prepTime); err == nil {
			req.PrepTime = val
		}
	}
	if cookTime := c.PostForm("cook_time"); cookTime != "" {
		if val, err := strconv.Atoi(cookTime); err == nil {
			req.CookTime = val
		}
	}
	if servings := c.PostForm("servings"); servings != "" {
		if val, err := strconv.Atoi(servings); err == nil {
			req.Servings = val
		}
	}

	if ingredientsJSON := c.PostForm("ingredients"); ingredientsJSON != "" {
		if err := json.Unmarshal([]byte(ingredientsJSON), &req.Ingredients); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ingredients format: " + err.Error()})
			return
		}
	}

	if stepsJSON := c.PostForm("steps"); stepsJSON != "" {
		if err := json.Unmarshal([]byte(stepsJSON), &req.Steps); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid steps format: " + err.Error()})
			return
		}
	}

	if req.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	thumbnail, _ := c.FormFile("thumbnail")

	res, err := h.Service.CreateRecipe(req, uid, thumbnail)
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
// @Success 200 {array} dto.RecipeResponse
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
// @Success 200 {object} dto.RecipeResponse
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
// @Success 200 {array} dto.RecipeResponse
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
// @Accept multipart/form-data
// @Produce json
// @Param thumbnail formData file false "New Thumbnail Image"
// @Success 200 {object} dto.RecipeResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/recipes/{id} [put]
// UpdateRecipe godoc
// @Summary Update recipe
// @Description Update recipe details by ID
// @Tags Recipes
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Recipe ID"
// @Param title formData string false "Recipe Title"
// @Param description formData string false "Recipe Description"
// @Param category formData string false "Recipe Category"
// @Param prep_time formData int false "Preparation Time"
// @Param cook_time formData int false "Cooking Time"
// @Param servings formData int false "Number of Servings"
// @Param ingredients formData string false "Ingredients JSON Array"
// @Param steps formData string false "Steps JSON Array"
// @Success 200 {object} dto.RecipeResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/recipes/{id} [put]
func (h *RecipeHandler) UpdateRecipe(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateRecipeRequest

	if title := c.PostForm("title"); title != "" {
		req.Title = &title
	}
	if description := c.PostForm("description"); description != "" {
		req.Description = &description
	}
	if category := c.PostForm("category"); category != "" {
		req.Category = &category
	}

	if prepTime := c.PostForm("prep_time"); prepTime != "" {
		if val, err := strconv.Atoi(prepTime); err == nil {
			req.PrepTime = &val
		}
	}
	if cookTime := c.PostForm("cook_time"); cookTime != "" {
		if val, err := strconv.Atoi(cookTime); err == nil {
			req.CookTime = &val
		}
	}
	if servings := c.PostForm("servings"); servings != "" {
		if val, err := strconv.Atoi(servings); err == nil {
			req.Servings = &val
		}
	}

	if ingredientsJSON := c.PostForm("ingredients"); ingredientsJSON != "" {
		var ingredients []dto.IngredientInput
		if err := json.Unmarshal([]byte(ingredientsJSON), &ingredients); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ingredients format: " + err.Error()})
			return
		}
		req.Ingredients = ingredients
	}

	if stepsJSON := c.PostForm("steps"); stepsJSON != "" {
		var steps []dto.StepInput
		if err := json.Unmarshal([]byte(stepsJSON), &steps); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid steps format: " + err.Error()})
			return
		}
		req.Steps = steps
	}

	thumbnail, _ := c.FormFile("thumbnail")

	res, err := h.Service.UpdateRecipe(id, req, thumbnail)
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
