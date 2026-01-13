package training

type LogTrainingRequest struct {
	ExerciseName string  `json:"exercise_name"`
	Weight       float64 `json:"weight"`
	Sets         int     `json:"sets"`
	Reps         int     `json:"reps"`
	RPE          float64 `json:"rpe"`
	Notes        string  `json:"notes"`
}

type LogTrainingResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}
