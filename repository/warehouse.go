package repository

import (
	"alfdwirhmn/inventory/database"
	"alfdwirhmn/inventory/model"
	"context"
	"database/sql"
	"errors"

	"go.uber.org/zap"
)

type WarehouseRepository interface {
	Create(ctx context.Context, whs *model.Warehouse) (*model.Warehouse, error)
	Lists(page, limit int) ([]model.Warehouse, int, error)
	DetailById(id int) (*model.Warehouse, error)
	Update(ctx context.Context, id int, payload *model.Warehouse) (*model.Warehouse, error)
	Delete(ctx context.Context, id int) error
}

type warehouseRepository struct {
	DB     database.PgxIface
	Logger *zap.Logger
}

func NewWarehouseRepository(db database.PgxIface, log *zap.Logger) WarehouseRepository {
	return &warehouseRepository{
		DB:     db,
		Logger: log,
	}
}

func (r *warehouseRepository) Create(ctx context.Context, whs *model.Warehouse) (*model.Warehouse, error) {
	query := `
		INSERT INTO warehouses ( code, name, address, city, province, postal_code, phone, is_active,
    	created_by) VALUES ($1,  $2,  $3,  $4,  $5,  $6,  $7, true, $8)
		RETURNING id, code, name, address, city, province, postal_code, phone, is_active, created_by, created_at, updated_at;
	`

	var wrhs model.Warehouse
	err := r.DB.QueryRow(ctx, query,
		whs.Code,
		whs.Name,
		whs.Address,
		whs.City,
		whs.Province,
		whs.PostalCode,
		whs.Phone,
		whs.CreatedBy,
	).Scan(
		&wrhs.ID,
		&wrhs.Code,
		&wrhs.Name,
		&wrhs.Address,
		&wrhs.City,
		&wrhs.Province,
		&wrhs.PostalCode,
		&wrhs.Phone,
		&wrhs.IsActive,
		&wrhs.CreatedBy,
		&wrhs.CreatedAt,
		&wrhs.UpdatedAt,
	)

	if err != nil {
		r.Logger.Error("failed create warehouse", zap.Error(err))
		return nil, err
	}

	r.Logger.Info("warehouse created successfully")
	return &wrhs, nil
}

func (r *warehouseRepository) Lists(page, limit int) ([]model.Warehouse, int, error) {

	offset := (page - 1) * limit

	var total int
	countQuery := `SELECT COUNT(*) FROM warehouses`
	if err := r.DB.QueryRow(context.Background(), countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
        SELECT
            id, code, name, address, city, province,
            postal_code, phone, is_active,
            created_by, created_at, updated_at
        FROM warehouses
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `

	rows, err := r.DB.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// warehouses := make([]model.Warehouse, 0)
	var warehouses []model.Warehouse

	for rows.Next() {
		var whs model.Warehouse
		if err := rows.Scan(
			&whs.ID,
			&whs.Code,
			&whs.Name,
			&whs.Address,
			&whs.City,
			&whs.Province,
			&whs.PostalCode,
			&whs.Phone,
			&whs.IsActive,
			&whs.CreatedBy,
			&whs.CreatedAt,
			&whs.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}

		warehouses = append(warehouses, whs)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return warehouses, total, nil
}

func (r *warehouseRepository) DetailById(id int) (*model.Warehouse, error) {
	query := `
		SELECT id, code, name, address, city, province, postal_code, phone, is_active, created_by, created_at, updated_at
		FROM warehouses WHERE id = $1
	`

	var wr model.Warehouse
	err := r.DB.QueryRow(context.Background(),
		query, id).Scan(
		&wr.ID,
		&wr.Code,
		&wr.Name,
		&wr.Address,
		&wr.City,
		&wr.Province,
		&wr.PostalCode,
		&wr.Phone,
		&wr.IsActive,
		&wr.CreatedBy,
		&wr.CreatedAt,
		&wr.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &wr, err
}

func (r *warehouseRepository) Update(ctx context.Context, id int, payload *model.Warehouse) (*model.Warehouse, error) {
	query := `
			UPDATE warehouses
			SET code = $1, name = $2, address = $3, city = $4, province = $5, postal_code = $6, phone = $7, is_active = $8, updated_at = CURRENT_TIMESTAMP
			WHERE id = $9
			RETURNING id, code, name, address, city, province, postal_code, phone, is_active, created_by, created_at, updated_at;
	`

	var warehouses model.Warehouse

	err := r.DB.QueryRow(ctx, query,
		payload.Code,
		payload.Name,
		payload.Address,
		payload.City,
		payload.Province,
		payload.PostalCode,
		payload.Phone,
		payload.IsActive,
		id,
	).Scan(
		&warehouses.ID,
		&warehouses.Code,
		&warehouses.Name,
		&warehouses.Address,
		&warehouses.City,
		&warehouses.Province,
		&warehouses.PostalCode,
		&warehouses.Phone,
		&warehouses.IsActive,
		&warehouses.CreatedBy,
		&warehouses.CreatedAt,
		&warehouses.UpdatedAt,
	)

	if err != nil {
		return nil, errors.New("warehouse not found or already deleted")
	}

	return &warehouses, nil
}

func (r *warehouseRepository) Delete(ctx context.Context, id int) error {
	query := `
		UPDATE warehouses
		SET
			is_active = false,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1;
	`

	result, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("warehouses not found or already deleted")
	}

	return nil
}
