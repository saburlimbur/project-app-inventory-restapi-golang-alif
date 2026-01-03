package service

import (
	"alfdwirhmn/inventory/dto"
	"alfdwirhmn/inventory/model"
	"alfdwirhmn/inventory/repository"
	"alfdwirhmn/inventory/utils"
	"context"
	"errors"
)

type CategoryService interface {
	Create(ctx context.Context, usr *model.User, req dto.CreateCategoryRequest) (*model.Category, error)
	FindAll(page, limit int) (*[]model.Category, *dto.Pagination, error)
}

type categoryService struct {
	repo    repository.CategoryRepository
	permSvc PermissionService
}

func NewCategoryService(repo repository.CategoryRepository, permSvc PermissionService) CategoryService {
	return &categoryService{
		repo:    repo,
		permSvc: permSvc,
	}
}

func (c *categoryService) Create(ctx context.Context, usr *model.User, req dto.CreateCategoryRequest) (*model.Category, error) {
	// cek permission role
	if !c.permSvc.CanCreateMasterData(usr.Role) {
		return nil, errors.New("forbidden: cannot create category")
	}

	// uniq name
	exists, err := c.repo.IsCategoryNameExists(ctx, req.Name, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("category name already exists")
	}

	// uniq code
	exists, err = c.repo.IsCategoryCodeExists(ctx, req.Code, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("category code already exists")
	}

	createdBy := usr.ID

	// mapping
	category := &model.Category{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
		CreatedBy:   &createdBy,
	}

	return c.repo.Create(ctx, category)
}

func (c *categoryService) FindAll(page, limit int) (*[]model.Category, *dto.Pagination, error) {
	category, total, err := c.repo.Lists(page, limit)

	if err != nil {
		return nil, nil, err
	}

	pagination := dto.Pagination{
		Page:       page,
		Limit:      limit,
		TotalPages: utils.TotalPage(limit, int64(total)),
		TotalRows:  total,
	}

	return &category, &pagination, nil
}
