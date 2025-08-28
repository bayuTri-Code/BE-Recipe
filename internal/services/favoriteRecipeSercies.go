package services

import (
	"fmt"

	"github.com/bayuTri-Code/BE-Recipe/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FavoriteService struct {
	DB *gorm.DB
}

func NewFavoriteService(db *gorm.DB) *FavoriteService {
	return &FavoriteService{DB: db}
}

func (s *FavoriteService) GetAllFavorites(userID string) ([]models.Favorite, error) {
	var favorites []models.Favorite

	
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	err = s.DB.
		Preload("Recipe", func(db *gorm.DB) *gorm.DB {
			return db.
				Preload("User").
				Preload("Ingredients").
				Preload("Steps", func(db *gorm.DB) *gorm.DB {
					return db.Order("steps.number ASC")
				}).
				Preload("Photos").
				Preload("Favorites")
		}).
		Where("user_id = ?", userUUID).
		Find(&favorites).Error

	if err != nil {
		return nil, err
	}
	return favorites, nil
}


func (s *FavoriteService) AddFavoriteService(userID, recipeID string) (bool, error) {
	recipeUUID, err := uuid.Parse(recipeID)
	if err != nil {
		return false, fmt.Errorf("invalid recipe ID")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return false, fmt.Errorf("invalid user ID")
	}

	var fav models.Favorite
	err = s.DB.Where("user_id = ? AND recipe_id = ?", userUUID, recipeUUID).First(&fav).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if err == gorm.ErrRecordNotFound {
		newFav := models.Favorite{
			ID:       uuid.New(),
			UserID:   userUUID,
			RecipeID: recipeUUID,
		}

		if err := s.DB.Create(&newFav).Error; err != nil {
			return false, err
		}
		return true, nil // added
	}

	
	if err := s.DB.Delete(&fav).Error; err != nil {
		return false, err
	}
	return false, nil // removed
}

func (s *FavoriteService) RemoveFavorite(userID, recipeID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}

	recipeUUID, err := uuid.Parse(recipeID)
	if err != nil {
		return fmt.Errorf("invalid recipe ID")
	}

	return s.DB.
		Where("user_id = ? AND recipe_id = ?", userUUID, recipeUUID).
		Delete(&models.Favorite{}).Error
}
