package models

import (
	"github.com/google/uuid"
)

type AuthData struct {
	ID       uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`
	Name     string    `gorm:"not null" json:"name"`
	Email    string    `gorm:"unique;not null" json:"email"`
	Password string    `gorm:"not null" json:"password"`
}
