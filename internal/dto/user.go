package dto

type LoginResponse struct {
	Token string `json:"token"`
}

type UserProfileResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email"`
	Password string `json:"password" validate:"required,min=6"`
}
