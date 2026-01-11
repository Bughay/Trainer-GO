package food

import "github.com/Bughay/Trainer-GO/db"

type CreateFoodItemRequest struct {
	FoodName    string  `json:"food_name"`
	Calories100 float64 `json:"calories_100"`
	Protein100  float64 `json:"protein_100"`
	Carbs100    float64 `json:"carbs_100"`
	Fats100     float64 `json:"fats_100"`
}

type CreateFoodItemResponse struct {
	Message string   `json:"message"`
	Success bool     `json:"success"`
	Food    FoodItem `json:"food,omitempty"`
}

type LogFoodItemRequest struct {
	FoodName   string  `json:"food_name"`
	TotalGrams float64 `json:"total_grams"`
	Calories   float64 `json:"calories"`
	Protein    float64 `json:"protein"`
	Carbs      float64 `json:"carbs"`
	Fats       float64 `json:"fats"`
}
type LogFoodItemResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type FoodItem struct {
	UserID      int64   `json:"user_id"`
	FoodName    string  `json:"food_name"`
	Calories100 float64 `json:"calories_100"`
	Protein100  float64 `json:"protein_100"`
	Carbs100    float64 `json:"carbs_100"`
	Fats100     float64 `json:"fats_100"`
}

type ViewFoodRequest struct {
	DateFrom string `json:"from"`
	DateTo   string `json:"to"`
}

type ViewFoodRow struct {
	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
	Carbs    float64 `json:"carbs"`
	Fats     float64 `json:"fats"`
}

type ViewFoodResponse struct {
	Message string           `json:"message"`
	Success bool             `json:"success"`
	Foods   []db.ViewFoodRow `json:"foods"`
}
