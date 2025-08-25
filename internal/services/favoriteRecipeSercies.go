package services

import (
	"errors"

	"github.com/bayuTri-Code/BE-Recipe/database"
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


func GetAllFavoritesService(userID string) ([]models.Favorite, error) {
	var favorites []models.Favorite
	if err := database.Db.Preload("Recipe").Where("user_id = ?", userID).Find(&favorites).Error; err != nil {
		return nil, err
	}
	return favorites, nil
}

func (s *FavoriteService) AddFavorite(recipeID string, userID uuid.UUID) error {
	var r models.Recipe
	if err := s.DB.First(&r, "id = ?", recipeID).Error; err != nil {
		return errors.New("recipe not found")
	}
	var u models.User
	if err := s.DB.First(&u, "id = ?", userID).Error; err != nil {
		return errors.New("user not found")
	}

	var existing models.Favorite
	if err := s.DB.First(&existing, "user_id = ? AND recipe_id = ?", userID, r.ID).Error; err == nil {
		return nil 
	}

	f := models.Favorite{
		ID:       uuid.New(),
		UserID:   userID,
		RecipeID: r.ID,
	}
	return s.DB.Create(&f).Error
}

func (s *FavoriteService) RemoveFavorite(recipeID string, userID uuid.UUID) error {
	return s.DB.
		Where("user_id = ? AND recipe_id = ?", userID, recipeID).
		Delete(&models.Favorite{}).Error
}
