package food

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
func timeToPgTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  t,
		Valid: true,
	}
}

func parseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		// Default to today
		return time.Now(), nil
	}

	// Try multiple date formats
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
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
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(CreateFoodItemResponse{
			Message: "Authentication required",
			Success: false,
		})
		return
	}
	params := db.CreateFoodItemParams{
		UserID:      userID,
		FoodName:    request.FoodName,
		Calories100: request.Calories100,
		Protein100:  request.Protein100,
		Carbs100:    request.Carbs100,
		Fats100:     request.Fats100,
	}
	foodItem, err := h.queries.CreateFoodItem(r.Context(), params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(CreateFoodItemResponse{
			Message: "failed to create food item",
			Success: false,
		})
		return
	}

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

func (h *FoodHandler) ViewFoodHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query()
	dateFromStr := query.Get("from")
	dateToStr := query.Get("to")
	if dateFromStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ViewFoodResponse{
			Message: "'from' date parameter is required. Format: YYYY-MM-DD",
			Success: false,
		})
		return
	}

	if dateToStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ViewFoodResponse{
			Message: "'to' date parameter is required. Format: YYYY-MM-DD",
			Success: false,
		})
		return
	}

	dateFrom, err := time.Parse("2006-01-02", dateFromStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ViewFoodResponse{
			Message: "Invalid 'from' date format. Use YYYY-MM-DD",
			Success: false,
		})
		return
	}

	dateTo, err := time.Parse("2006-01-02", dateToStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ViewFoodResponse{
			Message: "Invalid 'to' date format. Use YYYY-MM-DD",
			Success: false,
		})
		return
	}

	if dateTo.Before(dateFrom) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ViewFoodResponse{
			Message: "'to' date must be after 'from' date",
			Success: false,
		})
		return
	}

	userID, ok := r.Context().Value(auth.UserIDKey).(int64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(CreateFoodItemResponse{
			Message: "Authentication required",
			Success: false,
		})
		return
	}
	viewFoodParams := db.ViewFoodParams{
		UserID:      userID,
		CreatedAt:   timeToPgTimestamp(dateFrom),
		CreatedAt_2: timeToPgTimestamp(dateTo),
	}
	foods, err := h.queries.ViewFood(r.Context(), viewFoodParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ViewFoodResponse{
			Message: fmt.Sprintf("Failed to fetch food entries: %v", err),
			Success: false,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Println(foods)
	json.NewEncoder(w).Encode(ViewFoodResponse{
		Message: "Food entries retrieved successfully",
		Success: true,
		Foods:   foods,
	})

}

func (h *FoodHandler) ViewFoodTotalHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query()
	dateFromStr := query.Get("from")
	dateToStr := query.Get("to")

	if dateFromStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ViewFoodTotalResponse{
			Message: "'from' date parameter is required. Format: YYYY-MM-DD",
			Success: false,
		})
		return
	}

	if dateToStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ViewFoodTotalResponse{
			Message: "'to' date parameter is required. Format: YYYY-MM-DD",
			Success: false,
		})
		return
	}

	// Parse dates
	dateFrom, err := time.Parse("2006-01-02", dateFromStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ViewFoodTotalResponse{
			Message: "Invalid 'from' date format. Use YYYY-MM-DD",
			Success: false,
		})
		return
	}

	dateTo, err := time.Parse("2006-01-02", dateToStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ViewFoodTotalResponse{
			Message: "Invalid 'to' date format. Use YYYY-MM-DD",
			Success: false,
		})
		return
	}

	if dateTo.Before(dateFrom) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ViewFoodTotalResponse{
			Message: "'to' date must be after 'from' date",
			Success: false,
		})
		return
	}

	userID, ok := r.Context().Value(auth.UserIDKey).(int64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ViewFoodTotalResponse{
			Message: "Authentication required",
			Success: false,
		})
		return
	}

	viewFoodTotalParams := db.ViewFoodTotalParams{
		UserID:      userID,
		CreatedAt:   timeToPgTimestamp(dateFrom),
		CreatedAt_2: timeToPgTimestamp(dateTo),
	}

	totals, err := h.queries.ViewFoodTotal(r.Context(), viewFoodTotalParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ViewFoodTotalResponse{
			Message: fmt.Sprintf("Failed to fetch food totals: %v", err),
			Success: false,
		})
		return
	}

	response := ViewFoodTotalResponse{
		Message: fmt.Sprintf("Food totals from %s to %s",
			dateFrom.Format("2006-01-02"),
			dateTo.Format("2006-01-02")),
		Success: true,
		Totals: struct {
			Calories float64 `json:"calories"`
			Protein  float64 `json:"protein"`
			Carbs    float64 `json:"carbs"`
			Fats     float64 `json:"fats"`
		}{
			Calories: totals.TotalCalories,
			Protein:  totals.TotalProtein,
			Carbs:    totals.TotalCarbs,
			Fats:     totals.TotalFats,
		},
	}

	// Return successful response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
