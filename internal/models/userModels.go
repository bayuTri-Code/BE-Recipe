package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)
type User struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Bio       string         `gorm:"type:text" json:"bio"`
	Avatar    string         `gorm:"type:text" json:"avatar"`
	Banner   string         `gorm:"type:text" json:"banner"`

	Recipes   []Recipe   `gorm:"foreignKey:UserID" json:"recipes"`
	Favorites []Favorite `gorm:"foreignKey:UserID" json:"favorites"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
