package repository

import (
	"alfdwirhmn/inventory/database"
	"alfdwirhmn/inventory/model"
	"context"
	"errors"

	"go.uber.org/zap"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	Lists(page, limit int) ([]model.User, int, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int) error

	IsEmailExists(ctx context.Context, email string, excludeID int) (bool, error)
	IsUsernameExists(ctx context.Context, username string, excludeID int) (bool, error)
	FindByIdentifier(ctx context.Context, identifier string) (*model.User, error)
	FindByID(ctx context.Context, id int) (*model.User, error)
}

type userRepository struct {
	DB     database.PgxIface
	Logger *zap.Logger
}

func NewUserRepository(db database.PgxIface, log *zap.Logger) UserRepository {
	return &userRepository{
		DB:     db,
		Logger: log,
	}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	query := `
		INSERT INTO users (username, email, password_hash, full_name, role, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, username, email, password_hash, full_name, role, is_active, created_at, updated_at
	`

	var created model.User
	err := r.DB.QueryRow(ctx, query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.Role,
		user.IsActive,
	).Scan(
		&created.ID,
		&created.Username,
		&created.Email,
		&created.PasswordHash,
		&created.FullName,
		&created.Role,
		&created.IsActive,
		&created.CreatedAt,
		&created.UpdatedAt,
	)

	if err != nil {
		r.Logger.Error("Failed to create user", zap.Error(err))
		return nil, err
	}

	r.Logger.Info("User created successfully", zap.Int("user_id", created.ID))
	return &created, nil
}

func (r *userRepository) Lists(page, limit int) ([]model.User, int, error) {
	offset := (page - 1) * limit

	var total int
	// ambil semua role
	query := `
		SELECT id, username, email, full_name, role, is_active, created_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	err := r.DB.QueryRow(context.Background(), query).Scan(&total)

	rows, err := r.DB.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var usr model.User
		if err := rows.Scan(
			&usr.ID,
			&usr.Username,
			&usr.Email,
			&usr.FullName,
			&usr.Role,
			&usr.IsActive,
			&usr.CreatedAt,
		); err != nil {
			return nil, 0, err
		}

		users = append(users, usr)
	}

	return users, total, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users
		SET
			username = $1,
			email = $2,
			full_name = $3,
			role = $4,
			is_active = $5,
			updated_at = NOW()
		WHERE id = $6
	`

	cmd, err := r.DB.Exec(ctx, query,
		user.Username,
		user.Email,
		user.FullName,
		user.Role,
		user.IsActive,
		user.ID,
	)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	query := `
		UPDATE users
		SET is_active = false, updated_at = NOW()
		WHERE id = $1
	`

	cmd, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}

// uniq user field berdasarkan email
func (r *userRepository) IsEmailExists(ctx context.Context, email string, excludeID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND id != $2)`

	var exists bool
	err := r.DB.QueryRow(ctx, query, email, excludeID).Scan(&exists)
	if err != nil {
		r.Logger.Error("Failed to check email existence", zap.Error(err))
		return false, err
	}

	return exists, nil
}

// uniq user field berdasarkan username
func (r *userRepository) IsUsernameExists(ctx context.Context, username string, excludeID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 AND id != $2)`

	var exists bool
	err := r.DB.QueryRow(ctx, query, username, excludeID).Scan(&exists)
	if err != nil {
		r.Logger.Error("Failed to check username existence", zap.Error(err))
		return false, err
	}

	return exists, nil
}

// iedntifier field body, login berdasarkan email atau username
func (r *userRepository) FindByIdentifier(ctx context.Context, identifier string) (*model.User, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, role, is_active, created_at, updated_at
		FROM users
		WHERE email = $1 OR username = $1
	`

	var user model.User
	err := r.DB.QueryRow(ctx, query, identifier).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id int) (*model.User, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, role, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user model.User
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
