package dto

import "time"

type WarehouseResponseDTO struct {
	ID         int     `json:"id"`
	Code       string  `json:"code"`
	Name       string  `json:"name"`
	Address    *string `json:"address,omitempty"`
	City       *string `json:"city,omitempty"`
	Province   *string `json:"province,omitempty"`
	PostalCode *string `json:"postal_code,omitempty"`
	Phone      *string `json:"phone,omitempty"`
	IsActive   bool    `json:"is_active"`
	CreatedBy  int     `json:"created_by"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateWarehouseRequest struct {
	Code       string  `json:"code" validate:"required,max=20"`
	Name       string  `json:"name" validate:"required,max=100"`
	Address    *string `json:"address,omitempty"`
	City       *string `json:"city,omitempty"`
	Province   *string `json:"province,omitempty"`
	PostalCode *string `json:"postal_code,omitempty" validate:"max=10"`
	Phone      *string `json:"phone,omitempty" validate:"max=20"`
	// PostalCode *string `json:"postal_code,omitempty,max=10"`
	// Phone      *string `json:"phone,omitempty,max=20"`
}

type UpdateWarehouseRequest struct {
	Code       string  `json:"code" validate:"required,max=20"`
	Name       string  `json:"name" validate:"required,max=100"`
	Address    *string `json:"address,omitempty"`
	City       *string `json:"city,omitempty"`
	Province   *string `json:"province,omitempty"`
	PostalCode *string `json:"postal_code,omitempty" validate:"max=10"`
	Phone      *string `json:"phone,omitempty" validate:"max=20"`
	IsActive   bool    `json:"is_active,omitempty"`
}
