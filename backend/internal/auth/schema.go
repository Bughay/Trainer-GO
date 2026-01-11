package auth

type UserRegistrationRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRegistrationResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	UserID  string `json:"user_id,omitempty"`
}

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Token   string `json:"token,omitempty"`
}

type User struct {
	UserID         int64
	HashedPassword string
}
