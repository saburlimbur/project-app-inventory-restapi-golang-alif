package service

import (
	"alfdwirhmn/inventory/database"
	"alfdwirhmn/inventory/repository"

	"go.uber.org/zap"
)

type Container struct {
	User      UserService
	Auth      AuthService
	Category  CategoryService
	Warehouse WarehouseService
	Racks     RacksService
	Items     ItemsService
	Sale      SaleService
}

func NewContainer(repo *repository.Container, log *zap.Logger, tx database.TxManager) *Container {
	permSvc := NewPermissionService()

	return &Container{
		User:      NewUserService(repo.UserRepo, permSvc),
		Auth:      NewAuthService(repo.UserRepo, repo.SessionRepo),
		Category:  NewCategoryService(repo.CategoryRepo, permSvc),
		Warehouse: NewWarehouseService(repo.WarehouseRepo, permSvc),
		Racks:     NewRacksService(repo.RacksRepo, permSvc),
		Items:     NewItemsService(repo.ItemsRepo, permSvc),
		Sale: NewSaleService(
			repo.SaleRepo,
			repo.ItemsRepo,
			permSvc,
			tx,
			log,
		),
	}
}
