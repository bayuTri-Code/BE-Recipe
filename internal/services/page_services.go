package services

import (
	"errors"

	"github.com/bayuTri-Code/BE-Recipe/internal/models"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type DashboardService struct {
	DB *gorm.DB
}

func NewDashboardService(db *gorm.DB) *DashboardService {
	return &DashboardService{DB: db}
}

func (s *DashboardService) GetDashboardData(userID uuid.UUID) (map[string]interface{}, error) {
	var user models.User
	if err := s.DB.Preload("Recipes").Preload("Favorites").First(&user, "id = ?", userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	totalRecipes := len(user.Recipes)
	totalFavorites := len(user.Favorites)

	dashboard := map[string]interface{}{
		"username": user.Name,
		"stats": map[string]int{
			"numberOfRecipes": totalRecipes,
			"numberOfFavorites":    totalFavorites,
		},
	}

	return dashboard, nil
}
