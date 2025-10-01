package dto

type DashboardDTO struct {
	Stats    StatsDTO `json:"stats"`
	Username string   `json:"username"`
}

type StatsDTO struct {
	NumberOfFavorites    int `json:"numberOfFavorites"`
	NumberOfRecipes int `json:"numberOfRecipes"`
}
