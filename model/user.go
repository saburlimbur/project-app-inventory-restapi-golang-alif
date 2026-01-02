package model

import "time"

type User struct {
	ID           int    `json:"id" db:"id"`
	Username     string `json:"username" db:"username"`
	Email        string `json:"email" db:"email"`
	PasswordHash string `json:"-" db:"password_hash"` // tdk terlihat dari JSON response
	FullName     string `json:"full_name" db:"full_name"`
	Role         string `json:"role" db:"role"`
	IsActive     bool   `json:"is_active" db:"is_active"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
