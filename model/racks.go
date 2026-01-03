package model

import "time"

type Racks struct {
	ID          int     `db:"id"`
	WarehouseID int     `db:"warehouse_id"`
	Code        string  `db:"code"`
	Name        string  `db:"name"`
	Location    *string `db:"location"`
	Capacity    int     `db:"capacity"`
	Description *string `db:"description"`
	IsActive    bool    `db:"is_active"`
	CreatedBy   *int    `db:"created_by"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
