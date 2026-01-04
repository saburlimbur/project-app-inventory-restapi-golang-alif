package model

import "time"

type Item struct {
	ID int `db:"id"`

	CategoryID int  `db:"category_id"`
	RackID     *int `db:"rack_id"`

	SKU          string  `db:"sku"`
	Name         string  `db:"name"`
	Description  *string `db:"description"`
	Unit         string  `db:"unit"`
	Price        float64 `db:"price"`
	Cost         float64 `db:"cost"`
	Stock        int     `db:"stock"`
	MinimumStock int     `db:"minimum_stock"`
	Weight       float64 `db:"weight"`
	Dimensions   *string `db:"dimensions"`
	IsActive     bool    `db:"is_active"`
	CreatedBy    *int    `json:"created_by,omitempty" db:"created_by"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
