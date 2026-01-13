package training

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/Bughay/Trainer-GO/db"
	"github.com/Bughay/Trainer-GO/internal/auth"
	"github.com/jackc/pgx/v5/pgtype"
)

type TrainingHandler struct {
	queries *db.Queries
}

func NewTrainingHandler(q *db.Queries) *TrainingHandler {
	return &TrainingHandler{
		queries: q,
	}
}

func Float64ToNumeric(value float64) pgtype.Numeric {
	scaled := int64(value * 100)
	return pgtype.Numeric{
		Int:   big.NewInt(scaled),
		Exp:   -2,
		Valid: true,
	}
}

func IntToInt4(value int) pgtype.Int4 {
	if value == 0 {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: int32(value), Valid: true}
}

func StringToText(value string) pgtype.Text {
	if value == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: value, Valid: true}
}

func (h *TrainingHandler) LogTrainingHandler(w http.ResponseWriter, r *http.Request) {
	var request LogTrainingRequest
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(LogTrainingResponse{
			Message: "failed to decode request",
			Success: false,
		})
		return
	}

	userID, ok := r.Context().Value(auth.UserIDKey).(int64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(LogTrainingResponse{
			Message: "Authentication required",
			Success: false,
		})
		return
	}

	logExerciseParams := db.LogExerciseParams{
		UserID:       userID,
		ExerciseName: request.ExerciseName,
		Weight:       Float64ToNumeric(request.Weight),
		Sets:         int32(request.Sets),
		Reps:         int32(request.Reps),
		Rpe:          int32(request.RPE),
		Notes:        StringToText(request.Notes),
	}

	exerciseEntry, err := h.queries.LogExercise(r.Context(), logExerciseParams)

	if err != nil {
		fmt.Println("Error logging exercise:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(LogTrainingResponse{
			Message: "failed to log exercise",
			Success: false,
		})
		return
	}

	fmt.Println("Logged exercise:", exerciseEntry)

	response := LogTrainingResponse{
		Message: "success",
		Success: true,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
