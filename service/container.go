package service

import "alfdwirhmn/inventory/repository"

type Container struct {
	User     UserService
	Auth     AuthService
	Category CategoryService
}

func NewContainer(repo *repository.Container) *Container {
	permSvc := NewPermissionService()

	return &Container{
		User:     NewUserService(repo.UserRepo, permSvc),
		Auth:     NewAuthService(repo.UserRepo, repo.SessionRepo),
		Category: NewCategoryService(repo.CategoryRepo, permSvc),
	}
}
