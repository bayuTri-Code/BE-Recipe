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
		AllowOrigins:     []string{"http://localhost:5173", "http://192.168.100.8:3000", "http://recipe.com", "http://192.168.100.102:3000", "http://127.0.0.1:5173"},
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

	r.Static("/storage", "./public/storage")

	// Auth routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", handler.RegisterHandler)
		auth.POST("/login", handler.LoginHandler)
		auth.POST("/logout", middleware.AuthMiddleware(), handler.LogoutHandler)
	}

	// Recipe & Favorite routes
	recipeService := services.NewRecipeService(db)
	favoriteService := services.NewFavoriteService(db)

	recipeHandler := handler.NewRecipeHandler(recipeService)
	favoriteHandler := handler.NewFavoriteHandler(favoriteService)

	apiRecipe := r.Group("/api")
	{
		apiRecipe.GET("/recipes", recipeHandler.GetAllRecipes)

		apiRecipe.GET("/myrecipes", middleware.AuthMiddleware(), recipeHandler.GetMyRecipes)
		apiRecipe.GET("/recipes/:id", middleware.AuthMiddleware(), recipeHandler.GetRecipeByID)
		apiRecipe.POST("/recipes", middleware.AuthMiddleware(), recipeHandler.CreateRecipe)
		apiRecipe.PUT("/recipes/:id", middleware.AuthMiddleware(), recipeHandler.UpdateRecipe)
		apiRecipe.DELETE("/recipes/:id", middleware.AuthMiddleware(), recipeHandler.DeleteRecipe)


		// Favorites
		apiRecipe.GET("/recipes/favorites", middleware.AuthMiddleware(), favoriteHandler.GetAllFavorites)
		apiRecipe.POST("/recipes/:recipe_id/favorites", middleware.AuthMiddleware(), favoriteHandler.AddFavoriteHandler)
		// apiRecipe.DELETE("/recipes/:id/favorites/:user_id", favoriteHandler.RemoveFavorite)
	}

	// Dashboard routes
	dashboardService := services.NewDashboardService(db)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)

	apiDashboard := r.Group("/api/page")
	apiDashboard.Use(middleware.AuthMiddleware())
	{
		apiDashboard.GET("/dashboard", dashboardHandler.GetDashboard)
	}

	return r
}
