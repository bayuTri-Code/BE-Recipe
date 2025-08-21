package main

import (
    "fmt"

    "github.com/bayuTri-Code/Auth-Services/database"
    "github.com/bayuTri-Code/Auth-Services/internal/config"
    "github.com/bayuTri-Code/Auth-Services/internal/routes"

    _ "github.com/bayuTri-Code/Auth-Services/cmd/api/docs" 
    ginSwagger "github.com/swaggo/gin-swagger"
    swaggerFiles "github.com/swaggo/files"
)

// @title Auth Service API
// @version 1.0
// @description API documentation for Auth Service
// @host localhost:8080
// @BasePath /
func main() {
    config.ConfigDb()
    database.PostgresConn()

    // routes
    r := routes.Routes()

    // Swagger route
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    host := "0.0.0.0"
    port := "8080"

    fmt.Printf("server is running in http://%s:%s\n", host, port)
    r.Run(host + ":" + port)
}
