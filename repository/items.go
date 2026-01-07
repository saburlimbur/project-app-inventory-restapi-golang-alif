package repository

import (
	"alfdwirhmn/inventory/dto"
	"alfdwirhmn/inventory/model"
	"context"
	"errors"

	"go.uber.org/zap"
)

type ItemsRepository interface {
	Create(ctx context.Context, itm *model.Item) (*model.Item, error)
	Lists(page, limit int) ([]model.Item, int, error)
	Update(ctx context.Context, id int, req dto.UpdateItemRequest) (*model.Item, error)
	Delete(ctx context.Context, id int) error

	FindByID(ctx context.Context, id int) (*model.Item, error)
	ReduceStock(ctx context.Context, itemID int, qty int) error
}

type itemsRepository struct {
	DB DBTX
	// DB     database.PgxIface
	Logger *zap.Logger
}

func NewItemsRepository(db DBTX, log *zap.Logger) ItemsRepository {
	return &itemsRepository{
		DB:     db,
		Logger: log,
	}
}

func (r *itemsRepository) Create(ctx context.Context, itm *model.Item) (*model.Item, error) {
	query := `
		INSERT INTO items (
			category_id, rack_id, sku, name, description,
			unit, price, cost, stock, minimum_stock,
			weight, dimensions, is_active, created_by
		)
		VALUES (
			$1, $2, $3, $4, $5,
			COALESCE($6, 'pcs'),
			$7,
			COALESCE($8, 0),
			COALESCE($9, 0),
			COALESCE($10, 5),
			COALESCE($11, 0),
			$12,
			true,
			$13
		)
		RETURNING
			id, category_id, rack_id, sku, name, description,
			unit, price, cost, stock, minimum_stock,
			weight, dimensions, is_active,
			created_by, created_at, updated_at
	`

	var items model.Item
	err := r.DB.QueryRow(ctx, query,
		itm.CategoryID,
		itm.RackID,
		itm.SKU,
		itm.Name,
		itm.Description,
		itm.Unit,
		itm.Price,
		itm.Cost,
		itm.Stock,
		itm.MinimumStock,
		itm.Weight,
		itm.Dimensions,
		itm.CreatedBy,
	).Scan(
		&items.ID,
		&items.CategoryID,
		&items.RackID,
		&items.SKU,
		&items.Name,
		&items.Description,
		&items.Unit,
		&items.Price,
		&items.Cost,
		&items.Stock,
		&items.MinimumStock,
		&items.Weight,
		&items.Dimensions,
		&items.IsActive,
		&items.CreatedBy,
		&items.CreatedAt,
		&items.UpdatedAt,
	)

	if err != nil {
		r.Logger.Error("failed to create items", zap.Error(err))
		return nil, err
	}

	r.Logger.Info("items created succesfully")
	return &items, nil
}

func (r *itemsRepository) Lists(page, limit int) ([]model.Item, int, error) {
	offset := (page - 1) * limit

	var total int
	// item list semua yang aktif
	countQuery := `
		SELECT COUNT(*)
		FROM items
		WHERE is_active = true
	`
	if err := r.DB.QueryRow(context.Background(), countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT
			id, category_id, rack_id, sku, name, description,
			unit, price, cost, stock, minimum_stock,
			weight, dimensions, is_active,
			created_by, created_at, updated_at
		FROM items
		WHERE is_active = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.DB.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	var items []model.Item

	for rows.Next() {
		var itm model.Item
		if err := rows.Scan(
			&itm.ID,
			&itm.CategoryID,
			&itm.RackID,
			&itm.SKU,
			&itm.Name,
			&itm.Description,
			&itm.Unit,
			&itm.Price,
			&itm.Cost,
			&itm.Stock,
			&itm.MinimumStock,
			&itm.Weight,
			&itm.Dimensions,
			&itm.IsActive,
			&itm.CreatedBy,
			&itm.CreatedAt,
			&itm.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}

		items = append(items, itm)
	}

	return items, total, nil
}

func (r *itemsRepository) Update(ctx context.Context, id int, req dto.UpdateItemRequest) (*model.Item, error) {
	// query := `
	// 	UPDATE items
	// 	SET category_id = $1, rack_id = $2, sku = $3, name = $4, description = $5, unit = $6, price = $7, cost = $8, stock = $9, minimum_stock = $10, weight = $11, dimensions = $12, is_active = $13, updated_at = CURRENT_TIMESTAMP
	// 	WHERE id = $14
	// 	RETURNING id, category_id, rack_id, sku, name, description, unit, price, cost, stock, minimum_stock, weight, dimensions, is_active, created_by, created_at, updated_at
	// `
	query := `
	UPDATE items
	SET
		category_id   = COALESCE($1, category_id),
		rack_id       = COALESCE($2, rack_id),
		sku           = COALESCE($3, sku),
		name          = COALESCE($4, name),
		description   = COALESCE($5, description),
		unit          = COALESCE($6, unit),
		price         = COALESCE($7, price),
		cost          = COALESCE($8, cost),
		stock         = COALESCE($9, stock),
		minimum_stock = COALESCE($10, minimum_stock),
		weight        = COALESCE($11, weight),
		dimensions    = COALESCE($12, dimensions),
		is_active     = COALESCE($13, is_active),
		updated_at    = CURRENT_TIMESTAMP
	WHERE id = $14
	AND is_active = true
	RETURNING
			id, category_id, rack_id, sku, name, description,
			unit, price, cost, stock, minimum_stock,
			weight, dimensions, is_active,
			created_by, created_at, updated_at
	`

	var item model.Item

	err := r.DB.QueryRow(ctx, query,
		req.CategoryID,
		req.RackID,
		req.SKU,
		req.Name,
		req.Description,
		req.Unit,
		req.Price,
		req.Cost,
		req.Stock,
		req.MinimumStock,
		req.Weight,
		req.Dimensions,
		req.IsActive,
		id,
	).Scan(
		&item.ID,
		&item.CategoryID,
		&item.RackID,
		&item.SKU,
		&item.Name,
		&item.Description,
		&item.Unit,
		&item.Price,
		&item.Cost,
		&item.Stock,
		&item.MinimumStock,
		&item.Weight,
		&item.Dimensions,
		&item.IsActive,
		&item.CreatedBy,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		return nil, errors.New("item not found or already deleted")
	}

	return &item, nil
}

func (r *itemsRepository) Delete(ctx context.Context, id int) error {
	query := `
		UPDATE items
		SET is_active = false,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		  AND is_active = true
	`

	result, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("item not found or already delete")
	}

	return nil
}

func (r *itemsRepository) FindByID(ctx context.Context, id int) (*model.Item, error) {
	query := `SELECT id, category_id, rack_id, sku, name, description, price, unit, cost, stock, minimum_stock, weight, dimensions, is_active, created_by, created_at, updated_at
			FROM items 
			WHERE id = $1`

	item := &model.Item{}
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&item.ID,
		&item.CategoryID,
		&item.RackID,
		&item.SKU,
		&item.Name,
		&item.Description,
		&item.Price,
		&item.Unit,
		&item.Cost,
		&item.Stock,
		&item.MinimumStock,
		&item.Weight,
		&item.Dimensions,
		&item.IsActive,
		&item.CreatedBy,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (r *itemsRepository) ReduceStock(ctx context.Context, itemID, qty int) error {
	q := `
	UPDATE items
	SET stock = stock - $1
	WHERE id = $2 AND stock >= $1
	`

	res, err := r.DB.Exec(ctx, q, qty, itemID)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errors.New("stock not enough")
	}

	return nil
}
