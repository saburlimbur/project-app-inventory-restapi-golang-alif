package dto

import "time"

type SaleResponseDTO struct {
	ID            int    `json:"id"`
	InvoiceNumber string `json:"invoice_number"`

	CustomerName  *string   `json:"customer_name"`
	CustomerPhone *string   `json:"customer_phone"`
	CustomerEmail *string   `json:"customer_email"`
	SaleDate      time.Time `json:"sale_date"`

	TotalAmount float64 `json:"total_amount"`
	Discount    float64 `json:"discount"`
	Tax         float64 `json:"tax"`
	GrandTotal  float64 `json:"grand_total"`

	PaymentMethod *string `json:"payment_method"`
	PaymentStatus string  `json:"payment_status"`
	Notes         *string `json:"notes"`

	CreatedBy int       `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateSaleItemRequest struct {
	ItemID   int     `json:"item_id" validate:"required"`
	Quantity int     `json:"quantity" validate:"required,gt=0"`
	Discount float64 `json:"discount" validate:"gte=0"`
}

type CreateSaleRequest struct {
	InvoiceNumber string                  `json:"invoice_number" validate:"required"`
	CustomerName  string                  `json:"customer_name"`
	CustomerPhone string                  `json:"customer_phone"`
	CustomerEmail string                  `json:"customer_email"`
	Discount      float64                 `json:"discount" validate:"gte=0"`
	Tax           float64                 `json:"tax" validate:"gte=0"`
	PaymentMethod string                  `json:"payment_method"`
	Notes         string                  `json:"notes"`
	Items         []CreateSaleItemRequest `json:"items" validate:"required,min=1"`
}

type UpdateSaleRequest struct {
	CustomerName  *string `json:"customer_name"`
	CustomerPhone *string `json:"customer_phone"`
	CustomerEmail *string `json:"customer_email"`

	TotalAmount *float64 `json:"total_amount" validate:"omitempty,gte=0"`
	Discount    *float64 `json:"discount" validate:"omitempty,gte=0"`
	Tax         *float64 `json:"tax" validate:"omitempty,gte=0"`
	GrandTotal  *float64 `json:"grand_total" validate:"omitempty,gte=0"`

	PaymentMethod *string `json:"payment_method"`
	PaymentStatus *string `json:"payment_status" validate:"omitempty,oneof=pending paid cancelled"`
	Notes         *string `json:"notes"`
}
