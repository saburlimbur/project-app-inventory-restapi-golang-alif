package repository

import (
	"alfdwirhmn/inventory/database"

	"go.uber.org/zap"
)

type ItemsRepository interface {
}

type itemsRepository struct {
	DB     database.PgxIface
	Logger *zap.Logger
}

func NewItemsRepository(db database.PgxIface, log *zap.Logger) ItemsRepository {
	return &itemsRepository{
		DB:     db,
		Logger: log,
	}
}
