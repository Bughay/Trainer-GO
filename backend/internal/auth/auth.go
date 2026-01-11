package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Bughay/Trainer-GO/db"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	queries   *db.Queries
	jwtSecret []byte
}

func NewAuthHandler(q *db.Queries, jwtSecret string) (*AuthHandler, error) {
	if jwtSecret == "" {
		return nil, fmt.Errorf("jwt secret cannot be empty")
	}
	return &AuthHandler{
		queries:   q,
		jwtSecret: []byte(jwtSecret),
	}, nil
}

func (h *AuthHandler) UserRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	var request UserRegistrationRequest
	var response UserRegistrationResponse
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response = UserRegistrationResponse{
			Message: "failed",
			Success: false,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(UserRegistrationResponse{
			Message: "Failed to hash password",
			Success: false,
		})
		return
	}
	userParams := db.CreateUserParams{
		Username:       request.Username,
		HashedPassword: string(hashedPassword),
	}
	user, err := h.queries.CreateUser(r.Context(), userParams)

	if err != nil {
		// Other database errors
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(UserRegistrationResponse{
			Message: "Database error",
			Success: false,
		})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(UserRegistrationResponse{
		Message: "User registered successfully",
		Success: true,
		UserID:  fmt.Sprintf("%d", user.UserID),
	})
	return
}

func (h *AuthHandler) UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	var request UserLoginRequest
	var response UserLoginResponse
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response = UserLoginResponse{
			Message: "failed",
			Success: false,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := h.queries.GetUserByUsername(r.Context(), request.Username)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response = UserLoginResponse{
			Message: "Invalid username or password",
			Success: false,
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(request.Password))
	if err != nil {
		// Password doesn't match
		w.WriteHeader(http.StatusUnauthorized)
		response = UserLoginResponse{
			Message: "Invalid username or password", // Same message for security
			Success: false,
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	token, err := h.GenerateToken(int64(user.UserID), user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(UserLoginResponse{
			Message: "Failed to generate authentication token",
			Success: false,
		})
		return
	}

	response = UserLoginResponse{
		Message: "works", // Same message for security
		Success: true,
		Token:   token,
	}
	json.NewEncoder(w).Encode(response)

	return

}
