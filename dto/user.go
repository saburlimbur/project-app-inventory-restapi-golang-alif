package dto

import (
	"alfdwirhmn/inventory/model"
)

type UserResponseDTO struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// untuk admin dan super admin
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=8,max=100"`
	FullName string `json:"full_name" validate:"required,min=3,max=100"`
	Role     string `json:"role" validate:"required,oneof=super_admin admin staff"`
	IsActive *bool  `json:"is_active,omitempty"`
}

// publik
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=admin staff"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"` // email atau username
	Password   string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token     string           `json:"token"`
	ExpiredAt string           `json:"expired_at"`
	User      *UserResponseDTO `json:"user"`
}

type UpdateUserRequest struct {
	Username *string `json:"username" validate:"omitempty,min=3"`
	Email    *string `json:"email" validate:"omitempty,email"`
	FullName *string `json:"full_name" validate:"omitempty"`
	Role     *string `json:"role" validate:"omitempty,oneof=super_admin admin staff"`
	IsActive *bool   `json:"is_active"`
}

// helper konversi dari model user untuk response dto user
func ToUserResponseDTO(user *model.User) *UserResponseDTO {
	return &UserResponseDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
