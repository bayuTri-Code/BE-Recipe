package main

import (
	"fmt"

	_ "github.com/bayuTri-Code/BE-Recipe/cmd/api/docs"
	"github.com/bayuTri-Code/BE-Recipe/database"
	"github.com/bayuTri-Code/BE-Recipe/internal/config"
	"github.com/bayuTri-Code/BE-Recipe/internal/routes"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Auth Service API
// @version 1.0
// @description API documentation for Auth Service
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	config.ConfigDb()
	db := database.PostgresConn()


	r := routes.Routes(db)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	host := "0.0.0.0"
	port := "8080"

	fmt.Printf("server is running in http://%s:%s\n", host, port)
	r.Run(host + ":" + port)
}
