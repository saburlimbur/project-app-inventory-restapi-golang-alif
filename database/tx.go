package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TxManager interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

type txManager struct {
	DB *pgxpool.Pool
}

func NewTxManager(db *pgxpool.Pool) TxManager {
	return &txManager{DB: db}
}

func (t *txManager) Begin(ctx context.Context) (pgx.Tx, error) {
	return t.DB.Begin(ctx)
}
