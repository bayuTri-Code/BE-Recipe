package routes

import (
	"log"

	"github.com/bayuTri-Code/BE-Recipe/internal/handler"
	"github.com/bayuTri-Code/BE-Recipe/internal/middleware"
	"github.com/bayuTri-Code/BE-Recipe/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routes(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://192.168.100.8:3000", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	err := r.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		log.Panicf("Failed to set trusted proxies: %v", err)
	}

	r.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "test api",
		})
	})

	recipeService := services.NewRecipeService(db)
	favoriteService := services.NewFavoriteService(db)

	recipeHandler := handler.NewRecipeHandler(recipeService)
	favoriteHandler := handler.NewFavoriteHandler(favoriteService)

	auth := r.Group("/auth")
	{
		auth.POST("/register", handler.RegisterHandler)
		auth.POST("/login", handler.LoginHandler)
	}

	apiRecipe := r.Group("/api")
	apiRecipe.Use(middleware.AuthMiddleware())
	{
		apiRecipe.POST("/recipes", recipeHandler.CreateRecipe)
		apiRecipe.GET("/recipes", recipeHandler.GetAllRecipes)
		apiRecipe.GET("/recipes/:id", recipeHandler.GetRecipeByID)
		apiRecipe.PUT("/recipes/:id", recipeHandler.UpdateRecipe)
		apiRecipe.DELETE("/recipes/:id", recipeHandler.DeleteRecipe)


		apiRecipe.GET("/myrecipes", recipeHandler.GetMyRecipes)

		// Favorites
		apiRecipe.POST("/recipes/:recipe_id/favorites", favoriteHandler.AddFavoriteHandler)
		apiRecipe.DELETE("/recipes/:id/favorites/:user_id", favoriteHandler.RemoveFavorite)
		apiRecipe.GET("/recipes/favorites", favoriteHandler.GetAllFavorites)
	}
	return r
}
