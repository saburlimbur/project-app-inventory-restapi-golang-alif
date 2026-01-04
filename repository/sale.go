package repository

import (
	"alfdwirhmn/inventory/model"
	"context"

	"go.uber.org/zap"
)

type SaleRepository interface {
	Create(ctx context.Context, sl *model.Sale) (*model.Sale, error)
	Lists(ctx context.Context, page, limit int) ([]model.Sale, int, error)

	CreateItem(ctx context.Context, item *model.SaleItem) error
	FindByID(ctx context.Context, id int) (*model.Sale, error)
}

type saleRepository struct {
	DB     DBTX
	Logger *zap.Logger
}

func NewSaleRepository(db DBTX, log *zap.Logger) SaleRepository {
	return &saleRepository{
		DB:     db,
		Logger: log,
	}
}

func (s *saleRepository) Create(ctx context.Context, sl *model.Sale) (*model.Sale, error) {
	query := `
	INSERT INTO sales (invoice_number, customer_name, customer_phone, customer_email, total_amount, discount, tax, grand_total, payment_method, notes, created_by) 
	VALUES ($1, $2, $3, $4, $5, COALESCE($6, 0), COALESCE($7, 0), $8, $9, $10, $11)
	RETURNING id, invoice_number, customer_name, customer_phone, customer_email, sale_date, total_amount, discount, tax, grand_total, payment_method, payment_status, notes, created_by, created_at, updated_at;
	`

	var customerName, customerPhone, customerEmail, paymentMethod, notes *string
	var createdBy *int

	var sale model.Sale

	err := s.DB.QueryRow(ctx, query,
		sl.InvoiceNumber,
		sl.CustomerName,
		sl.CustomerPhone,
		sl.CustomerEmail,
		sl.TotalAmount,
		sl.Discount,
		sl.Tax,
		sl.GrandTotal,
		sl.PaymentMethod,
		sl.Notes,
		sl.CreatedBy,
	).Scan(
		&sale.ID,
		&sale.InvoiceNumber,
		&customerName,
		&customerPhone,
		&customerEmail,
		&sale.SaleDate,
		&sale.TotalAmount,
		&sale.Discount,
		&sale.Tax,
		&sale.GrandTotal,
		&paymentMethod,
		&sale.PaymentStatus,
		&notes,
		&createdBy,
		&sale.CreatedAt,
		&sale.UpdatedAt,
	)

	if err != nil {
		s.Logger.Error("failed to create sale", zap.Error(err))
		return nil, err
	}

	// assign variabel pointer sementara ke struct
	sale.CustomerName = customerName
	sale.CustomerPhone = customerPhone
	sale.CustomerEmail = customerEmail
	sale.PaymentMethod = paymentMethod
	sale.Notes = notes
	sale.CreatedBy = createdBy

	s.Logger.Info("sale created successfully", zap.Int("id", sale.ID))
	return &sale, nil
}

func (s *saleRepository) Lists(ctx context.Context, page, limit int) ([]model.Sale, int, error) {
	offset := (page - 1) * limit

	var total int
	countQuery := `
		SELECT COUNT(*) FROM sales;
	`
	if err := s.DB.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
			SELECT
			    id,
			    invoice_number,
			    customer_name,
			    sale_date,
				total_amount,
				discount,
				tax,
			    grand_total,
			    payment_status,
			    payment_method,
			    created_at
			FROM sales
			ORDER BY sale_date DESC
			LIMIT $1 OFFSET $2;

	`
	rows, err := s.DB.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	var sale []model.Sale

	for rows.Next() {
		var sl model.Sale
		if err := rows.Scan(
			&sl.ID,
			&sl.InvoiceNumber,
			&sl.CustomerName,
			&sl.SaleDate,
			&sl.TotalAmount,
			&sl.Discount,
			&sl.Tax,
			&sl.GrandTotal,
			&sl.PaymentStatus,
			&sl.PaymentMethod,
			&sl.CreatedAt,
		); err != nil {
			return nil, 0, err
		}

		sale = append(sale, sl)
	}

	return sale, total, nil
}

func (s *saleRepository) CreateItem(ctx context.Context, item *model.SaleItem) error {
	query := `
	INSERT INTO sale_items (sale_id, item_id, quantity, unit_price, subtotal, discount)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id;
	`

	var id int
	err := s.DB.QueryRow(ctx, query,
		item.SaleID,
		item.ItemID,
		item.Quantity,
		item.UnitPrice,
		item.Subtotal,
		item.Discount,
	).Scan(&id)

	if err != nil {
		s.Logger.Error("failed to create sale item", zap.Error(err))
		return err
	}

	item.ID = id
	s.Logger.Info("sale item created", zap.Int("id", id))
	return nil
}

func (r *saleRepository) FindByID(ctx context.Context, id int) (*model.Sale, error) {
	q := `SELECT id, invoice_number, grand_total FROM sales WHERE id=$1`

	s := &model.Sale{}
	err := r.DB.QueryRow(ctx, q, id).Scan(
		&s.ID,
		&s.InvoiceNumber,
		&s.GrandTotal,
	)
	if err != nil {
		return nil, err
	}
	return s, nil
}
