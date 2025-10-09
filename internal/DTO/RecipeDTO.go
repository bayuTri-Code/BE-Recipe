package dto

import "github.com/google/uuid"

// ---- Requests ----
type CreateRecipeRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Category    string `json:"category"`
	PrepTime    int    `json:"prep_time"`
	CookTime    int    `json:"cook_time"`
	Servings    int    `json:"servings"`

	Ingredients []IngredientInput `json:"ingredients"`
	Steps       []StepInput       `json:"steps"`
}

type UpdateRecipeRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
	PrepTime    *int    `json:"prep_time"`
	CookTime    *int    `json:"cook_time"`
	Servings    *int    `json:"servings"`

	Ingredients []IngredientInput `json:"ingredients"`
	Steps       []StepInput       `json:"steps"`
}

type IngredientInput struct {
	Name   string `json:"name" binding:"required"`
	Amount string `json:"amount" binding:"required"`
}

type StepInput struct {
	Number int    `json:"number" binding:"required"`
	Detail string `json:"detail" binding:"required"`
}

type UserSummaryResponse struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
	Avatar string    `json:"avatar"`
}

type IngredientResponse struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Amount string    `json:"amount"`
}

type StepResponse struct {
	ID     uuid.UUID `json:"id"`
	Number int       `json:"number"`
	Detail string    `json:"detail"`
}

type FavoriteResponse struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

type RecipeResponse struct {
	ID          uuid.UUID             `json:"id"`
	Title       string                `json:"title"`
	Description string                `json:"description"`
	Category    string                `json:"category"`
	Thumbnail   string                `json:"thumbnail"`
	User        UserSummaryResponse `json:"user"`
	Ingredients []IngredientResponse  `json:"ingredients"`
	Steps       []StepResponse        `json:"steps"`
	PrepTime    int                   `json:"prep_time"`
	CookTime    int                   `json:"cook_time"`
	Servings    int                   `json:"servings"`
	Favorites   []FavoriteResponse    `json:"favorites"`
}

type AddFavoriteRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
}
