package dto

import "time"

type RackResponseDTO struct {
	ID          int       `json:"id"`
	WarehouseID int       `json:"warehouse_id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Location    *string   `json:"location,omitempty"`
	Capacity    int       `json:"capacity"`
	Description *string   `json:"description,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedBy   int       `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateRackRequest struct {
	WarehouseID int     `json:"warehouse_id" validate:"required"`
	Code        string  `json:"code" validate:"required,max=20"`
	Name        string  `json:"name" validate:"required,max=100"`
	Location    *string `json:"location,omitempty"`
	Capacity    int     `json:"capacity,omitempty"`
	Description *string `json:"description,omitempty"`
}

type UpdateRackRequest struct {
	Code        string  `json:"code" validate:"required,max=20"`
	Name        string  `json:"name" validate:"required,max=100"`
	Location    *string `json:"location,omitempty"`
	Capacity    int     `json:"capacity,omitempty"`
	Description *string `json:"description,omitempty"`
	IsActive    bool    `json:"is_active"`
}
