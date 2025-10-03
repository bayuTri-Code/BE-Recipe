package database

import (
	"fmt"
	"log"

	dto "github.com/bayuTri-Code/BE-Recipe/internal/DTO"
	"github.com/bayuTri-Code/BE-Recipe/internal/config"
	"github.com/bayuTri-Code/BE-Recipe/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func PostgresConn() *gorm.DB {
	configDb := config.DbConfig

	SetDb := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		configDb.DBHost,
		configDb.DBPort,
		configDb.DBUser,
		configDb.DBPassword,
		configDb.DBName,
		configDb.DBSslmode,
	)

	db, err := gorm.Open(postgres.Open(SetDb), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	Db = db
	log.Println("Database Connected")

	autoMigrate(db)

	return db
}

func autoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Recipe{},
		&models.Ingredient{},
		&models.Step{},
		&models.Favorite{},
		&dto.BlacklistedToken{},
	)

	if err != nil {
		log.Fatalf("Auto Migration Failed: %v", err)
	}
	log.Println("Auto Migration Complete!")
}