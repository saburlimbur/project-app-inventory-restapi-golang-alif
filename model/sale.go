package model

import "time"

type Sale struct {
	ID            int
	InvoiceNumber string
	CustomerName  *string
	CustomerPhone *string
	CustomerEmail *string
	SaleDate      time.Time
	TotalAmount   float64
	Discount      float64
	Tax           float64
	GrandTotal    float64
	PaymentMethod *string
	PaymentStatus string
	Notes         *string
	CreatedBy     *int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type SaleItem struct {
	ID        int
	SaleID    int
	ItemID    int
	Quantity  int
	UnitPrice float64
	Subtotal  float64
	Discount  float64
	CreatedAt time.Time
}

// type Sale struct {
// 	ID            int       `db:"id"`
// 	InvoiceNumber string    `db:"invoice_number"`
// 	CustomerName  *string   `db:"customer_name"`
// 	CustomerPhone *string   `db:"customer_phone"`
// 	CustomerEmail *string   `db:"customer_email"`
// 	SaleDate      time.Time `db:"sale_date"`

// 	TotalAmount float64 `db:"total_amount"`
// 	Discount    float64 `db:"discount"`
// 	Tax         float64 `db:"tax"`
// 	GrandTotal  float64 `db:"grand_total"`

// 	PaymentMethod *string `db:"payment_method"`
// 	PaymentStatus string  `db:"payment_status"`
// 	Notes         *string `db:"notes"`

// 	CreatedBy *int      `json:"created_by,omitempty" db:"created_by"`
// 	CreatedAt time.Time `db:"created_at"`
// 	UpdatedAt time.Time `db:"updated_at"`
// }
