package services

import (
	"errors"

	"github.com/bayuTri-Code/BE-Recipe/database"
	dto "github.com/bayuTri-Code/BE-Recipe/internal/DTO"
	"github.com/bayuTri-Code/BE-Recipe/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RecipeService struct {
	DB *gorm.DB
}

func NewRecipeService(db *gorm.DB) *RecipeService {
	return &RecipeService{DB: db}
}

func toUserSummary(u models.User) dto.UserSummaryResponse {
	return dto.UserSummaryResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}

func toIngredientResponses(items []models.Ingredient) []dto.IngredientResponse {
	out := make([]dto.IngredientResponse, 0, len(items))
	for _, it := range items {
		out = append(out, dto.IngredientResponse{
			ID:     it.ID,
			Name:   it.Name,
			Amount: it.Amount,
		})
	}
	return out
}

func toStepResponses(items []models.Step) []dto.StepResponse {
	out := make([]dto.StepResponse, 0, len(items))
	for _, it := range items {
		out = append(out, dto.StepResponse{
			ID:     it.ID,
			Number: it.Number,
			Detail: it.Detail,
		})
	}
	return out
}

func toPhotoResponses(items []models.Photo) []dto.PhotoResponse {
	out := make([]dto.PhotoResponse, 0, len(items))
	for _, it := range items {
		out = append(out, dto.PhotoResponse{
			ID:  it.ID,
			URL: it.URL,
		})
	}
	return out
}

func toFavoriteResponses(items []models.Favorite) []dto.FavoriteResponse {
	out := make([]dto.FavoriteResponse, 0, len(items))
	for _, it := range items {
		out = append(out, dto.FavoriteResponse{
			ID:     it.ID,
			UserID: it.UserID,
		})
	}
	return out
}

func toRecipeResponse(m models.Recipe) dto.RecipeResponse {
	return dto.RecipeResponse{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		Category:    m.Category,
		User:        toUserSummary(m.User),
		Ingredients: toIngredientResponses(m.Ingredients),
		Steps:       toStepResponses(m.Steps),
		PrepTime:    m.PrepTime,
		CookTime:    m.CookTime,
		Servings:    m.Servings,
		Photos:      toPhotoResponses(m.Photos),
		Favorites:   toFavoriteResponses(m.Favorites),
	}
}

//crud
func (s *RecipeService) CreateRecipe(req dto.CreateRecipeRequest) (dto.RecipeResponse, error) {
	var out dto.RecipeResponse

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// ensure user exists (optional but recommended)
		var user models.User
		if err := tx.First(&user, "id = ?", req.UserID).Error; err != nil {
			return errors.New("user not found")
		}

		recipe := models.Recipe{
			ID:          uuid.New(),
			Title:       req.Title,
			Description: req.Description,
			Category:    req.Category,
			PrepTime:    req.PrepTime,
			CookTime:    req.CookTime,
			Servings:    req.Servings,
			UserID:      req.UserID,
		}
		if err := tx.Create(&recipe).Error; err != nil {
			return err
		}

		if len(req.Ingredients) > 0 {
			ings := make([]models.Ingredient, 0, len(req.Ingredients))
			for _, in := range req.Ingredients {
				ings = append(ings, models.Ingredient{
					ID:       uuid.New(),
					RecipeID: recipe.ID,
					Name:     in.Name,
					Amount:   in.Amount,
				})
			}
			if err := tx.Create(&ings).Error; err != nil {
				return err
			}
			recipe.Ingredients = ings
		}

		if len(req.Steps) > 0 {
			steps := make([]models.Step, 0, len(req.Steps))
			for _, st := range req.Steps {
				steps = append(steps, models.Step{
					ID:       uuid.New(),
					RecipeID: recipe.ID,
					Number:   st.Number,
					Detail:   st.Detail,
				})
			}
			if err := tx.Create(&steps).Error; err != nil {
				return err
			}
			recipe.Steps = steps
		}

		if len(req.Photos) > 0 {
			photos := make([]models.Photo, 0, len(req.Photos))
			for _, ph := range req.Photos {
				photos = append(photos, models.Photo{
					ID:       uuid.New(),
					RecipeID: recipe.ID,
					URL:      ph.URL,
				})
			}
			if err := tx.Create(&photos).Error; err != nil {
				return err
			}
			recipe.Photos = photos
		}

		recipe.User = user
		out = toRecipeResponse(recipe)
		return nil
	})

	return out, err
}

func (s *RecipeService) GetAllRecipes() ([]dto.RecipeResponse, error) {
	var list []models.Recipe
	err := s.DB.
		Preload("User").
		Preload("Ingredients").
		Preload("Steps", func(db *gorm.DB) *gorm.DB {
			return db.Order("steps.number ASC")
		}).
		Preload("Photos").
		Preload("Favorites").
		Find(&list).Error
	if err != nil {
		return nil, err
	}

	out := make([]dto.RecipeResponse, 0, len(list))
	for _, r := range list {
		out = append(out, toRecipeResponse(r))
	}
	return out, nil
}

func (s *RecipeService) GetRecipeByID(id string) (dto.RecipeResponse, error) {
	var r models.Recipe
	err := s.DB.
		Preload("User").
		Preload("Ingredients").
		Preload("Steps", func(db *gorm.DB) *gorm.DB {
			return db.Order("steps.number ASC")
		}).
		Preload("Photos").
		Preload("Favorites").
		First(&r, "id = ?", id).Error
	if err != nil {
		return dto.RecipeResponse{}, err
	}
	return toRecipeResponse(r), nil
}

func (s *RecipeService) UpdateRecipe(id string, req dto.UpdateRecipeRequest) (dto.RecipeResponse, error) {
	var out dto.RecipeResponse
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		var r models.Recipe
		if err := tx.First(&r, "id = ?", id).Error; err != nil {
			return errors.New("recipe not found")
		}

		if req.Title != nil {
			r.Title = *req.Title
		}
		if req.Description != nil {
			r.Description = *req.Description
		}
		if req.Category != nil {
			r.Category = *req.Category
		}
		if req.PrepTime != nil {
			r.PrepTime = *req.PrepTime
		}
		if req.CookTime != nil {
			r.CookTime = *req.CookTime
		}
		if req.Servings != nil {
			r.Servings = *req.Servings
		}

		if err := tx.Save(&r).Error; err != nil {
			return err
		}

		if req.Ingredients != nil {
			if err := tx.Where("recipe_id = ?", r.ID).Delete(&models.Ingredient{}).Error; err != nil {
				return err
			}
			if len(req.Ingredients) > 0 {
				ings := make([]models.Ingredient, 0, len(req.Ingredients))
				for _, in := range req.Ingredients {
					ings = append(ings, models.Ingredient{
						ID:       uuid.New(),
						RecipeID: r.ID,
						Name:     in.Name,
						Amount:   in.Amount,
					})
				}
				if err := tx.Create(&ings).Error; err != nil {
					return err
				}
			}
		}

		if req.Steps != nil {
			if err := tx.Where("recipe_id = ?", r.ID).Delete(&models.Step{}).Error; err != nil {
				return err
			}
			if len(req.Steps) > 0 {
				steps := make([]models.Step, 0, len(req.Steps))
				for _, st := range req.Steps {
					steps = append(steps, models.Step{
						ID:       uuid.New(),
						RecipeID: r.ID,
						Number:   st.Number,
						Detail:   st.Detail,
					})
				}
				if err := tx.Create(&steps).Error; err != nil {
					return err
				}
			}
		}

		if req.Photos != nil {
			if err := tx.Where("recipe_id = ?", r.ID).Delete(&models.Photo{}).Error; err != nil {
				return err
			}
			if len(req.Photos) > 0 {
				photos := make([]models.Photo, 0, len(req.Photos))
				for _, ph := range req.Photos {
					photos = append(photos, models.Photo{
						ID:       uuid.New(),
						RecipeID: r.ID,
						URL:      ph.URL,
					})
				}
				if err := tx.Create(&photos).Error; err != nil {
					return err
				}
			}
		}

		if err := tx.
			Preload("User").
			Preload("Ingredients").
			Preload("Steps", func(db *gorm.DB) *gorm.DB { return db.Order("steps.number ASC") }).
			Preload("Photos").
			Preload("Favorites").
			First(&r, "id = ?", r.ID).Error; err != nil {
			return err
		}

		out = toRecipeResponse(r)
		return nil
	})
	return out, err
}
func DeleteRecipeService(id string) error {
    var recipe models.Recipe

    if err := database.Db.First(&recipe, "id = ?", id).Error; err != nil {
        return err
    }

    if err := database.Db.Unscoped().Delete(&recipe).Error; err != nil {
        return err
    }

    return nil
}
