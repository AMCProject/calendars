package internal

const (
	Normal    = "standard"
	Semanal   = "weekly"
	Ocasional = "occasional"
)

type Calendar struct {
	UserId string `db:"user_id" json:"user_id"`
	MealId string `db:"meal_id" json:"meal_id"`
	Name   string `json:"name" json:"name"`
	Date   string `db:"date" json:"date"`
}

//definitions for endpoint calls//

type User struct {
	Id       string  `db:"id" json:"id,omitempty"`
	Name     *string `db:"name" json:"name,omitempty"`
	Mail     string  `db:"mail" json:"mail" validate:"required,excludes= "`
	Password string  `db:"password" json:"password" validate:"required,excludes= "`
}

type Meal struct {
	Id          string `db:"id" json:"id,omitempty"`
	UserId      string `db:"user_id" json:"userId"`
	Name        string `db:"name" json:"name" validate:"required"`
	Description string `db:"description" json:"description,omitempty"`
	Image       string `db:"image" json:"image,omitempty"`
	Type        string `db:"type" json:"type" validate:"required,oneof=weekly occasional standard"`
	Ingredients string `db:"ingredients" json:"ingredients" validate:"required"`
	Kcal        int    `db:"kcal" json:"kcal"`
	Seasons     string `db:"seasons" json:"seasons"`
	//Creator     int    `db:"creator" json:"creator"`
	//Saves       int    `db:"saves" json:"saves"`
}

type MealToFront struct {
	Id          string   `json:"id"`
	UserId      string   `json:"userId"`
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description"`
	Image       string   `json:"image"`
	Type        string   `json:"type" validate:"required,oneof=weekly occasional standard"`
	Ingredients []string `json:"ingredients" validate:"required"`
	Kcal        int      `json:"kcal"`
	Seasons     []string `json:"seasons" validate:"required,dive,oneof=summer winter spring fall general"`
	//Creator     int      `json:"creator"`
	//Saves       int      `json:"saves"`
}
