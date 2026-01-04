package repository

import (
	"alfdwirhmn/inventory/database"

	"go.uber.org/zap"
)

type Container struct {
	UserRepo      UserRepository
	CategoryRepo  CategoryRepository
	WarehouseRepo WarehouseRepository
	RacksRepo     RacksRepository
	ItemsRepo     ItemsRepository

	SessionRepo SessionRepository
}

func NewContainer(db database.PgxIface, log *zap.Logger) *Container {
	return &Container{
		UserRepo:      NewUserRepository(db, log),
		CategoryRepo:  NewCategoryRepository(db, log),
		WarehouseRepo: NewWarehouseRepository(db, log),
		RacksRepo:     NewRacksRepository(db, log),
		ItemsRepo:     NewItemsRepository(db, log),

		SessionRepo: NewSessionRepository(db),
		// SessionRepo: NewSessionRepository(db, log),
	}
}
