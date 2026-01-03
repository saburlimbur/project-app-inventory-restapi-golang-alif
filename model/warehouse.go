package model

import "time"

type Warehouse struct {
	ID         int     `db:"id"`
	Code       string  `db:"code"`
	Name       string  `db:"name"`
	Address    *string `db:"address"`
	City       *string `db:"city"`
	Province   *string `db:"province"`
	PostalCode *string `db:"postal_code"`
	Phone      *string `db:"phone"`
	IsActive   bool    `db:"is_active"`
	CreatedBy  *int    `db:"created_by"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
