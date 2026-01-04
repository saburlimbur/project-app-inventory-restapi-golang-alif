package handler

import (
	"alfdwirhmn/inventory/repository"
	"alfdwirhmn/inventory/service"
	"alfdwirhmn/inventory/utils"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Container struct {
	User      *UserHandler
	Auth      *AuthHandler
	Category  *CategoryHandler
	Warehouse *WarehouseHandler
	Racks     *RacksHandler
	Items     *ItemsHandler

	Repositories *repository.Container
}

func NewContainer(
	svc *service.Container,
	repo *repository.Container,
	log *zap.Logger,
	conf utils.Configuration,
) *Container {

	validate := validator.New()

	return &Container{
		User:      NewUserHandler(svc.User, validate, log, conf),
		Auth:      NewAuthHandler(svc.Auth, log),
		Category:  NewCategoryHandler(svc.Category, validate, log, conf),
		Warehouse: NewWarehouseHandler(svc.Warehouse, validate, log, conf),
		Racks:     NewRacksHandler(svc.Racks, validate, log, conf),
		Items:     NewItemsHandler(svc.Items, validate, log, conf),

		Repositories: repo,
	}
}
