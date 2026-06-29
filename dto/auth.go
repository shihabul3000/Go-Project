package dto

import "time"

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=120"`
	Email    string `json:"email" validate:"required,email,max=160"`
	Password string `json:"password" validate:"required,min=8,max=72"`
	Role     string `json:"role" validate:"omitempty,oneof=driver admin"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=160"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuthUserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type LoginResponse struct {
	Token string           `json:"token"`
	User  AuthUserResponse `json:"user"`
}
