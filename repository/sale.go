package repository

import (
	"alfdwirhmn/inventory/database"

	"go.uber.org/zap"
)

type SaleRepository interface {
}

type saleRepository struct {
	DB     database.PgxIface
	Logger *zap.Logger
}

func NewSaleRepository(db database.PgxIface, log *zap.Logger) SaleRepository {
	return &saleRepository{
		DB:     db,
		Logger: log,
	}
}
