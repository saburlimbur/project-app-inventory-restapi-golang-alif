package repository

import (
	"alfdwirhmn/inventory/database"
	"alfdwirhmn/inventory/model"
	"context"

	"go.uber.org/zap"
)

type CategoryRepository interface {
	Create(ctx context.Context, ctg *model.Category) (*model.Category, error)
	Lists(page, limit int) ([]model.Category, int, error)

	IsCategoryNameExists(ctx context.Context, name string, id int) (bool, error)
	IsCategoryCodeExists(ctx context.Context, code string, id int) (bool, error)
}

type categoryRepository struct {
	DB     database.PgxIface
	Logger *zap.Logger
}

func NewCategoryRepository(db database.PgxIface, log *zap.Logger) CategoryRepository {
	return &categoryRepository{
		DB:     db,
		Logger: log,
	}
}

func (r *categoryRepository) Create(ctx context.Context, ctg *model.Category) (*model.Category, error) {
	query := `
		INSERT INTO categories (code, name, description, is_active, created_by)
		VALUES ($1, $2, $3, true, $4)
		RETURNING id, code, name, description, is_active, created_by, created_at, updated_at
	`

	var ctgr model.Category
	err := r.DB.QueryRow(ctx, query,
		ctg.Code,
		ctg.Name,
		ctg.Description,
		ctg.CreatedBy,
	).Scan(
		&ctgr.ID,
		&ctgr.Code,
		&ctgr.Name,
		&ctgr.Description,
		&ctgr.IsActive,
		&ctgr.CreatedBy,
		&ctgr.CreatedAt,
		&ctgr.UpdatedAt,
	)

	if err != nil {
		r.Logger.Error("failed to create category", zap.Error(err))
		return nil, err
	}

	r.Logger.Info("category created succesfully")
	return &ctgr, nil
}

func (r *categoryRepository) Lists(page, limit int) ([]model.Category, int, error) {
	offset := (page - 1) * limit

	var totalCtg int

	query := `
			SELECT id, code, name, description, is_active, created_by, created_at
			FROM categories
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
	`

	err := r.DB.QueryRow(context.Background(), query).Scan(&totalCtg)

	rows, err := r.DB.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var category []model.Category

	for rows.Next() {
		var ctg model.Category
		if err := rows.Scan(
			&ctg.ID,
			&ctg.Code,
			&ctg.Name,
			&ctg.Description,
			&ctg.IsActive,
			&ctg.CreatedBy,
			&ctg.CreatedAt,
		); err != nil {
			return nil, 0, err
		}

		category = append(category, ctg)
	}

	return category, totalCtg, nil
}

func (r *categoryRepository) IsCategoryNameExists(
	ctx context.Context,
	name string,
	id int,
) (bool, error) {

	query := `
		SELECT EXISTS (SELECT 1 FROM categories WHERE LOWER(name) = LOWER($1) AND id != $2)
		`

	var exists bool
	err := r.DB.QueryRow(ctx, query, name, id).Scan(&exists)
	if err != nil {
		r.Logger.Error("failed to check category name exists", zap.Error(err))
		return false, err
	}

	return exists, nil
}

func (r *categoryRepository) IsCategoryCodeExists(ctx context.Context, code string, id int) (bool, error) {

	query := `SELECT EXISTS (SELECT 1 FROM categories WHERE code = $1 AND id != $2)
			`

	var exists bool
	err := r.DB.QueryRow(ctx, query, code, id).Scan(&exists)
	if err != nil {
		r.Logger.Error("failed to check category code exists", zap.Error(err))
		return false, err
	}

	return exists, nil
}
