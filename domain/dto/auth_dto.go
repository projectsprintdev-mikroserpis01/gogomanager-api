package dto

import "github.com/google/uuid"

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=100,ascii"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100,ascii"`
}

type RegisterResponse struct {
	ID uuid.UUID `json:"id"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100,ascii"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}
