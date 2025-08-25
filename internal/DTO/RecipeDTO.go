package dto

import "github.com/google/uuid"

// ---- Requests ----
type CreateRecipeRequest struct {
	Title       string             `json:"title" binding:"required"`
	Description string             `json:"description"`
	Category    string             `json:"category"`
	PrepTime    int                `json:"prep_time"`
	CookTime    int                `json:"cook_time"`
	Servings    int                `json:"servings"`
	UserID      uuid.UUID          `json:"user_id" binding:"required"`
	Ingredients []IngredientInput  `json:"ingredients"` // optional, bisa kosong
	Steps       []StepInput        `json:"steps"`       // optional
	Photos      []PhotoInput       `json:"photos"`      // optional
}

type UpdateRecipeRequest struct {
	Title       *string            `json:"title"`
	Description *string            `json:"description"`
	Category    *string            `json:"category"`
	PrepTime    *int               `json:"prep_time"`
	CookTime    *int               `json:"cook_time"`
	Servings    *int               `json:"servings"`
	// Prinsip replace: bila array dikirim, kita hapus semua lama & buat ulang baru
	Ingredients []IngredientInput  `json:"ingredients"`
	Steps       []StepInput        `json:"steps"`
	Photos      []PhotoInput       `json:"photos"`
}

type IngredientInput struct {
	Name   string `json:"name" binding:"required"`
	Amount string `json:"amount" binding:"required"`
}

type StepInput struct {
	Number int    `json:"number" binding:"required"`
	Detail string `json:"detail" binding:"required"`
}

type PhotoInput struct {
	URL string `json:"url" binding:"required"`
}

// ---- Responses ----
type UserSummaryResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
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

type PhotoResponse struct {
	ID  uuid.UUID `json:"id"`
	URL string    `json:"url"`
}

type FavoriteResponse struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

type RecipeResponse struct {
	ID          uuid.UUID            `json:"id"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Category    string               `json:"category"`
	User        UserSummaryResponse  `json:"user"`
	Ingredients []IngredientResponse `json:"ingredients"`
	Steps       []StepResponse       `json:"steps"`
	PrepTime    int                  `json:"prep_time"`
	CookTime    int                  `json:"cook_time"`
	Servings    int                  `json:"servings"`
	Photos      []PhotoResponse      `json:"photos"`
	Favorites   []FavoriteResponse   `json:"favorites"`
}


type AddFavoriteRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
}