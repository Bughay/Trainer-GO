package food

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Bughay/Trainer-GO/db"
	"github.com/Bughay/Trainer-GO/internal/auth"
	"github.com/jackc/pgx/v5/pgtype"
)

type FoodHandler struct {
	queries *db.Queries
}

func NewFoodHandler(q *db.Queries) *FoodHandler {
	return &FoodHandler{
		queries: q,
	}
}
func int64ToPgInt8(value int64, valid bool) pgtype.Int8 {
	return pgtype.Int8{
		Int64: value,
		Valid: valid,
	}
}

func (h *FoodHandler) CreateFoodItemHandler(w http.ResponseWriter, r *http.Request) {
	var request CreateFoodItemRequest
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CreateFoodItemResponse{
			Message: "failed to decode request",
			Success: false,
		})
		return
	}
	userID, ok := r.Context().Value(auth.UserIDKey).(int64)
	if !ok {
		// This means auth middleware wasn't run or failed
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(CreateFoodItemResponse{
			Message: "Authentication required",
			Success: false,
		})
		return
	}
	params := db.CreateFoodItemParams{
		UserID:      userID, // ⚠️ But this should come from JWT!
		FoodName:    request.FoodName,
		Calories100: request.Calories100,
		Protein100:  request.Protein100,
		Carbs100:    request.Carbs100,
		Fats100:     request.Fats100,
	}
	// 4. Call database
	foodItem, err := h.queries.CreateFoodItem(r.Context(), params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(CreateFoodItemResponse{
			Message: "failed to create food item",
			Success: false,
		})
		return
	}

	// 5. Return response
	response := CreateFoodItemResponse{
		Message: "Food item created",
		Success: true,
		Food: FoodItem{
			UserID:      foodItem.UserID,
			FoodName:    foodItem.FoodName,
			Calories100: foodItem.Calories100,
			Protein100:  foodItem.Protein100,
			Carbs100:    foodItem.Carbs100,
			Fats100:     foodItem.Fats100,
		},
	}
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(response)
}

func (h *FoodHandler) LogFoodHandler(w http.ResponseWriter, r *http.Request) {
	var request LogFoodItemRequest
	// var response LogFoodItemResponse
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CreateFoodItemResponse{
			Message: "failed to decode request",
			Success: false,
		})
		return
	}
	userID, ok := r.Context().Value(auth.UserIDKey).(int64)
	if !ok {
		// This means auth middleware wasn't run or failed
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(CreateFoodItemResponse{
			Message: "Authentication required",
			Success: false,
		})
		return
	}
	foodCacheParams := db.CreateFoodCacheItemParams{
		UserID:      userID,
		FoodName:    request.FoodName,
		Calories100: (request.Calories / request.TotalGrams),
		Protein100:  (request.Protein / request.TotalGrams),
		Carbs100:    (request.Carbs / request.TotalGrams),
		Fats100:     (request.Fats / request.TotalGrams),
	}
	logfoodCache, err := h.queries.CreateFoodCacheItem(r.Context(), foodCacheParams)

	logFoodParams := db.LogFoodItemParams{
		UserID:     userID,
		FoodID:     int64ToPgInt8(logfoodCache.FoodID, true),
		RecipeID:   int64ToPgInt8(0, false),
		Calories:   request.Calories,
		TotalGrams: request.TotalGrams,
		Protein:    request.Protein,
		Carbs:      request.Carbs,
		Fats:       request.Fats,
	}
	logfood, err := h.queries.LogFoodItem(r.Context(), logFoodParams)
	fmt.Println(logfood)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(CreateFoodItemResponse{
			Message: "failed to log food item",
			Success: false,
		})
		return
	}

	response := LogFoodItemResponse{
		Message: "success",
		Success: true,
	}
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(response)
}
