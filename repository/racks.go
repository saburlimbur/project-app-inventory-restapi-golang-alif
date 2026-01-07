package repository

import (
	"alfdwirhmn/inventory/database"
	"alfdwirhmn/inventory/model"
	"context"
	"database/sql"
	"errors"

	"go.uber.org/zap"
)

type RacksRepository interface {
	Create(ctx context.Context, rk *model.Racks) (*model.Racks, error)
	Lists(page, limit int) ([]model.Racks, int, error)
	DetailById(id int) (*model.Racks, error)
	Update(ctx context.Context, id int, payload *model.Racks) (*model.Racks, error)
	Delete(ctx context.Context, id int) error
}

type racksRepository struct {
	DB     database.PgxIface
	Logger *zap.Logger
}

func NewRacksRepository(db database.PgxIface, log *zap.Logger) RacksRepository {
	return &racksRepository{
		DB:     db,
		Logger: log,
	}
}

func (r *racksRepository) Create(ctx context.Context, rk *model.Racks) (*model.Racks, error) {
	query := `
		INSERT INTO racks (warehouse_id, code, name, location, capacity, description, is_active, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, true, $7)
		RETURNING id, warehouse_id, code, name, location, capacity, description, is_active, created_by, created_at, updated_at;
	`

	var rack model.Racks
	err := r.DB.QueryRow(ctx, query,
		rk.WarehouseID,
		rk.Code,
		rk.Name,
		rk.Location,
		rk.Capacity,
		rk.Description,
		rk.CreatedBy,
	).Scan(
		&rack.ID,
		&rack.WarehouseID,
		&rack.Code,
		&rack.Name,
		&rack.Location,
		&rack.Capacity,
		&rack.Description,
		&rack.IsActive,
		&rack.CreatedBy,
		&rack.CreatedAt,
		&rack.UpdatedAt,
	)

	if err != nil {
		r.Logger.Error("failed create rack", zap.Error(err))
		return nil, err
	}

	r.Logger.Info("rack created successfully")
	return &rack, nil
}

func (r *racksRepository) Lists(page, limit int) ([]model.Racks, int, error) {
	offset := (page - 1) * limit

	var total int
	countQuery := `SELECT COUNT(*) FROM racks WHERE is_active = true;`
	if err := r.DB.QueryRow(context.Background(), countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
	SELECT
		id, warehouse_id, code, name, location,
		capacity, description, is_active,
		created_by, created_at, updated_at
	FROM racks
	WHERE is_active = true
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2;
	`

	rows, err := r.DB.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var racks []model.Racks
	for rows.Next() {
		var rk model.Racks
		if err := rows.Scan(
			&rk.ID,
			&rk.WarehouseID,
			&rk.Code,
			&rk.Name,
			&rk.Location,
			&rk.Capacity,
			&rk.Description,
			&rk.IsActive,
			&rk.CreatedBy,
			&rk.CreatedAt,
			&rk.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		racks = append(racks, rk)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return racks, total, nil
}

func (r *racksRepository) DetailById(id int) (*model.Racks, error) {
	query := `
		SELECT id, warehouse_id, code, name, location, capacity, description, is_active, created_by, created_at, updated_at
		FROM racks
		WHERE id = $1
	`

	var rack model.Racks

	err := r.DB.QueryRow(context.Background(),
		query, id).Scan(
		&rack.ID,
		&rack.WarehouseID,
		&rack.Code,
		&rack.Name,
		&rack.Location,
		&rack.Capacity,
		&rack.Description,
		&rack.IsActive,
		&rack.CreatedBy,
		&rack.CreatedAt,
		&rack.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &rack, err
}

func (r *racksRepository) Update(ctx context.Context, id int, payload *model.Racks) (*model.Racks, error) {
	query := `
	UPDATE racks
	SET
		code = $1,
		name = $2,
		location = $3,
		capacity = $4,
		description = $5,
		is_active = $6,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $7
	RETURNING
		id, warehouse_id, code, name, location, capacity, description, is_active,
		created_by, created_at, updated_at;
	`

	var rk model.Racks
	err := r.DB.QueryRow(ctx, query,
		payload.Code,
		payload.Name,
		payload.Location,
		payload.Capacity,
		payload.Description,
		payload.IsActive,
		id,
	).Scan(
		&rk.ID,
		&rk.WarehouseID,
		&rk.Code,
		&rk.Name,
		&rk.Location,
		&rk.Capacity,
		&rk.Description,
		&rk.IsActive,
		&rk.CreatedBy,
		&rk.CreatedAt,
		&rk.UpdatedAt,
	)

	if err != nil {
		return nil, errors.New("rack not found or already deleted")
	}

	return &rk, nil
}

func (r *racksRepository) Delete(ctx context.Context, id int) error {
	query := `
	UPDATE racks
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
		return errors.New("rack not found or already deleted")
	}

	return nil
}
