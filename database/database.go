package database

import (
	"alfdwirhmn/inventory/utils"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxIface interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
}

func InitDB(config utils.DatabaseCofig) (*pgxpool.Pool, error) {
	fmt.Println("called config:", config)

	return nil, nil
}
