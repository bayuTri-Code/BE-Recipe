package services

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

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
		Avatar: u.Avatar,
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
		Thumbnail:   m.Thumbnail,
		User:        toUserSummary(m.User),
		Ingredients: toIngredientResponses(m.Ingredients),
		Steps:       toStepResponses(m.Steps),
		PrepTime:    m.PrepTime,
		CookTime:    m.CookTime,
		Servings:    m.Servings,
		Favorites:   toFavoriteResponses(m.Favorites),
	}
}

func (s *RecipeService) SaveThumbnail(file *multipart.FileHeader) (string, error) {
	apiImagePath := os.Getenv("API_IMAGE_PATH")

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
	if !allowedExts[ext] {
		return "", errors.New("only jpg, jpeg, and png files are allowed")
	}

	if file.Size > 2*1024*1024 {
		return "", errors.New("file size must not exceed 2MB")
	}

	filename := fmt.Sprintf("recipe_%d%s", time.Now().UnixNano(), ext)
	storagePath := "public/storage"
	fullPath := filepath.Join(storagePath, filename)

	if err := os.MkdirAll(storagePath, os.ModePerm); err != nil {
		return "", err
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/storage/%s", apiImagePath, filename), nil
}

func (s *RecipeService) DeleteThumbnail(thumbnailURL string) error {
	if thumbnailURL == "" {
		return nil
	}

	parts := strings.Split(thumbnailURL, "/storage/")
	if len(parts) != 2 {
		return nil
	}

	filename := parts[1]
	fullPath := filepath.Join("public/storage", filename)

	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

func (s *RecipeService) CreateRecipe(req dto.CreateRecipeRequest, userID uuid.UUID, thumbnail *multipart.FileHeader) (dto.RecipeResponse, error) {
	var out dto.RecipeResponse

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		var user models.User
		if err := tx.First(&user, "id = ?", userID).Error; err != nil {
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
			UserID:      userID,
		}

		if thumbnail != nil {
			thumbnailURL, err := s.SaveThumbnail(thumbnail)
			if err != nil {
				return fmt.Errorf("failed to save thumbnail: %w", err)
			}
			recipe.Thumbnail = thumbnailURL
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
		Preload("Favorites").
		First(&r, "id = ?", id).Error
	if err != nil {
		return dto.RecipeResponse{}, err
	}
	return toRecipeResponse(r), nil
}

func (s *RecipeService) UpdateRecipe(id string, req dto.UpdateRecipeRequest, thumbnail *multipart.FileHeader) (dto.RecipeResponse, error) {
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

		if thumbnail != nil {
			if err := s.DeleteThumbnail(r.Thumbnail); err != nil {
				return fmt.Errorf("failed to delete old thumbnail: %w", err)
			}

			thumbnailURL, err := s.SaveThumbnail(thumbnail)
			if err != nil {
				return fmt.Errorf("failed to save new thumbnail: %w", err)
			}
			r.Thumbnail = thumbnailURL
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

		if err := tx.
			Preload("User").
			Preload("Ingredients").
			Preload("Steps", func(db *gorm.DB) *gorm.DB { return db.Order("steps.number ASC") }).
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

	if recipe.Thumbnail != "" {
		parts := strings.Split(recipe.Thumbnail, "/storage/")
		if len(parts) == 2 {
			filename := parts[1]
			fullPath := filepath.Join("public/storage", filename)
			os.Remove(fullPath)
		}
	}

	if err := database.Db.Unscoped().Delete(&recipe).Error; err != nil {
		return err
	}

	return nil
}

func (s *RecipeService) GetRecipesByUserID(userID string) ([]dto.RecipeResponse, error) {
	var recipes []models.Recipe

	err := s.DB.
		Preload("User").
		Preload("Ingredients").
		Preload("Steps", func(db *gorm.DB) *gorm.DB {
			return db.Order("steps.number ASC")
		}).
		Preload("Favorites").
		Where("user_id = ?", userID).
		Find(&recipes).Error
	if err != nil {
		return nil, err
	}

	out := make([]dto.RecipeResponse, 0, len(recipes))
	for _, r := range recipes {
		out = append(out, toRecipeResponse(r))
	}

	return out, nil
}
