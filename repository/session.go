package repository

import (
	"alfdwirhmn/inventory/database"
	"alfdwirhmn/inventory/model"
	"context"
	"time"

	"github.com/google/uuid"
)

type SessionRepository interface {
	Create(ctx context.Context, s *model.Session) error
	FindByToken(ctx context.Context, token uuid.UUID) (*model.Session, error)
	Revoke(ctx context.Context, token uuid.UUID) error
}

type sessionRepository struct {
	DB database.PgxIface
}

func NewSessionRepository(db database.PgxIface) SessionRepository {
	return &sessionRepository{DB: db}
}

func (r *sessionRepository) Create(ctx context.Context, s *model.Session) error {
	query := `
		INSERT INTO sessions (user_id, token, expired_at, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.DB.Exec(ctx, query,
		s.UserID,
		s.Token,
		s.ExpiredAt,
		s.IPAddress,
		s.UserAgent,
	)
	return err
}

func (r *sessionRepository) FindByToken(ctx context.Context, token uuid.UUID) (*model.Session, error) {
	query := `
		SELECT id, user_id, token, created_at, expired_at, revoked_at
		FROM sessions
		WHERE token = $1
	`

	var s model.Session
	err := r.DB.QueryRow(ctx, query, token).Scan(
		&s.ID,
		&s.UserID,
		&s.Token,
		&s.CreatedAt,
		&s.ExpiredAt,
		&s.RevokedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *sessionRepository) Revoke(ctx context.Context, token uuid.UUID) error {
	query := `
		UPDATE sessions
		SET revoked_at = $1
		WHERE token = $2
	`
	_, err := r.DB.Exec(ctx, query, time.Now(), token)
	return err
}
