package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)



func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

type Recipe struct {
	ID          uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Category    string         `gorm:"type:varchar(50)" json:"category"`
	PrepTime    int            `json:"prep_time"`
	CookTime    int            `json:"cook_time"`
	Servings    int            `json:"servings"`
	UserID      uuid.UUID      `gorm:"type:char(36);index" json:"user_id"`

	// Relations
	User        User         `gorm:"foreignKey:UserID" json:"user"`
	Ingredients []Ingredient `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"ingredients"`
	Steps       []Step       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"steps"`
	Photos      []Photo      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"photos"`
	Favorites   []Favorite   `gorm:"foreignKey:RecipeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"favorites"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (r *Recipe) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return
}

type Ingredient struct {
	ID       uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	RecipeID uuid.UUID `gorm:"type:char(36);index" json:"recipe_id"`
	Name     string    `gorm:"not null" json:"name"`
	Amount   string    `gorm:"not null" json:"amount"`
}

func (i *Ingredient) BeforeCreate(tx *gorm.DB) (err error) {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return
}

type Step struct {
	ID       uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	RecipeID uuid.UUID `gorm:"type:char(36);index" json:"recipe_id"`
	Number   int       `gorm:"not null" json:"number"`
	Detail   string    `gorm:"type:text;not null" json:"detail"`
}

func (s *Step) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}

type Photo struct {
	ID       uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	RecipeID uuid.UUID `gorm:"type:char(36);index" json:"recipe_id"`
	URL      string    `gorm:"not null" json:"url"`
}

func (p *Photo) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}

type Favorite struct {
	ID       uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	UserID   uuid.UUID `gorm:"type:char(36);index" json:"user_id"`
	RecipeID uuid.UUID `gorm:"type:char(36);index" json:"recipe_id"`

	
	_ struct{} `gorm:"uniqueIndex:uniq_fav_user_recipe,priority:1"`
	
}

func (f *Favorite) BeforeCreate(tx *gorm.DB) (err error) {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return
}
