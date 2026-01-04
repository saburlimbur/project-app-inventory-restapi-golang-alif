package dto

import "time"

type ItemResponseDTO struct {
	ID int `json:"id"`

	CategoryID int  `json:"category_id"`
	RackID     *int `json:"rack_id,omitempty"`

	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`

	Unit         string  `json:"unit"`
	Price        float64 `json:"price"`
	Cost         float64 `json:"cost"`
	Stock        int     `json:"stock"`
	MinimumStock int     `json:"minimum_stock"`

	Weight     float64 `json:"weight"`
	Dimensions *string `json:"dimensions,omitempty"`
	IsActive   bool    `json:"is_active"`
	CreatedBy  int     `json:"created_by"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateItemRequest struct {
	CategoryID int  `json:"category_id" validate:"required"`
	RackID     *int `json:"rack_id,omitempty"`

	SKU         string  `json:"sku" validate:"required,max=50"`
	Name        string  `json:"name" validate:"required,max=200"`
	Description *string `json:"description,omitempty"`

	Unit         string  `json:"unit" validate:"omitempty,max=20"`
	Price        float64 `json:"price" validate:"required,gte=0"`
	Cost         float64 `json:"cost" validate:"omitempty,gte=0"`
	Stock        int     `json:"stock" validate:"omitempty,gte=0"`
	MinimumStock int     `json:"minimum_stock" validate:"omitempty,gte=0"`

	Weight     float64 `json:"weight" validate:"omitempty,gte=0"`
	Dimensions *string `json:"dimensions,omitempty"`
	IsActive   *bool   `json:"is_active"`
}

type UpdateItemRequest struct {
	CategoryID   *int     `json:"category_id" validate:"omitempty"`
	RackID       *int     `json:"rack_id" validate:"omitempty"`
	SKU          *string  `json:"sku" validate:"omitempty"`
	Name         *string  `json:"name" validate:"omitempty"`
	Description  *string  `json:"description" validate:"omitempty"`
	Unit         *string  `json:"unit" validate:"omitempty"`
	Price        *float64 `json:"price" validate:"omitempty,gt=0"`
	Cost         *float64 `json:"cost" validate:"omitempty,gte=0"`
	Stock        *int     `json:"stock" validate:"omitempty,gte=0"`
	MinimumStock *int     `json:"minimum_stock" validate:"omitempty,gte=0"`
	Weight       *float64 `json:"weight" validate:"omitempty,gte=0"`
	Dimensions   *string  `json:"dimensions" validate:"omitempty"`
	IsActive     *bool    `json:"is_active" validate:"omitempty"`
}
