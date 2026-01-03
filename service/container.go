package service

import "alfdwirhmn/inventory/repository"

type Container struct {
	User      UserService
	Auth      AuthService
	Category  CategoryService
	Warehouse WarehouseService
}

func NewContainer(repo *repository.Container) *Container {
	permSvc := NewPermissionService()

	return &Container{
		User:      NewUserService(repo.UserRepo, permSvc),
		Auth:      NewAuthService(repo.UserRepo, repo.SessionRepo),
		Category:  NewCategoryService(repo.CategoryRepo, permSvc),
		Warehouse: NewWarehouseService(repo.WarehouseRepo, permSvc),
	}
}
