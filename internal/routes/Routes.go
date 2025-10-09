package routes

import (
	"log"
	"os"
	"strings"

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
	originallow := strings.Split(os.Getenv("CORS_ORIGINS"), ",")
	r.Use(cors.New(cors.Config{
		AllowOrigins:     originallow,
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
	r.Static("/profile-storage", "./public/profile_storage")
	r.Static("/profile-banner", "./public/profile_banner")

	// Auth routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", middleware.RateLimiter(5, 60), handler.RegisterHandler)
		auth.POST("/login", middleware.RateLimiter(5, 60), handler.LoginHandler)
		auth.POST("/forgot-password", handler.ForgotPasswordHandler)
		auth.POST("/reset-password", handler.ResetPasswordHandler)
		auth.POST("/logout", middleware.AuthMiddleware(), handler.LogoutHandler)
		auth.PUT("/profile", middleware.AuthMiddleware(), middleware.RateLimiter(5, 60), handler.UpdateProfileHandler)
	}

	// Recipe & Favorite routes
	recipeService := services.NewRecipeService(db)
	favoriteService := services.NewFavoriteService(db)

	recipeHandler := handler.NewRecipeHandler(recipeService)
	favoriteHandler := handler.NewFavoriteHandler(favoriteService)

	apiRecipe := r.Group("/api")
	{
		apiRecipe.GET("/recipes", recipeHandler.GetAllRecipes)
		apiRecipe.GET("/recipesByCategory", handler.GetRecipesByCategory)

		apiRecipe.GET("/myrecipes", middleware.AuthMiddleware(), recipeHandler.GetMyRecipes)
		apiRecipe.GET("/recipes/:id", middleware.AuthMiddleware(), recipeHandler.GetRecipeByID)
		apiRecipe.POST("/recipes", middleware.AuthMiddleware(), middleware.RateLimiter(5, 60), recipeHandler.CreateRecipe)
		apiRecipe.PUT("/recipes/:id", middleware.AuthMiddleware(),  middleware.RateLimiter(10, 60), recipeHandler.UpdateRecipe)
		apiRecipe.DELETE("/recipes/:id", middleware.AuthMiddleware(), middleware.RateLimiter(15, 60), recipeHandler.DeleteRecipe)

		// Favorites
		apiRecipe.GET("/recipes/favorites", middleware.AuthMiddleware(), favoriteHandler.GetAllFavorites)
		apiRecipe.POST("/recipes/:recipe_id/favorites", middleware.AuthMiddleware(), middleware.RateLimiter(20, 60), favoriteHandler.AddFavoriteHandler)
		// apiRecipe.DELETE("/recipes/:id/favorites/:user_id", favoriteHandler.RemoveFavorite)
	}

	// Dashboard routes
	dashboardService := services.NewDashboardService(db)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)

	apiDashboard := r.Group("/api/page")
	apiDashboard.Use(middleware.AuthMiddleware(), middleware.RateLimiter(100, 60))
	{
		apiDashboard.GET("/dashboard", dashboardHandler.GetDashboard)
	}

	return r
}
