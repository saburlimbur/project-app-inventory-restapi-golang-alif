package service

import (
	"alfdwirhmn/inventory/dto"
	"alfdwirhmn/inventory/model"
	"alfdwirhmn/inventory/repository"
	"alfdwirhmn/inventory/utils"
	"context"
	"errors"
)

type ItemsService interface {
	Create(ctx context.Context, usr *model.User, req dto.CreateItemRequest) (*model.Item, error)
	FindAll(page, limit int) (*[]model.Item, *dto.Pagination, error)
	Update(ctx context.Context, usr *model.User, id int, req dto.UpdateItemRequest) (*model.Item, error)
	Delete(ctx context.Context, usr *model.User, id int) error
	FindByID(ctx context.Context, id int, usr *model.User) (*model.Item, error)
}

type itemsService struct {
	repo    repository.ItemsRepository
	permSvc PermissionService
}

func NewItemsService(repo repository.ItemsRepository, permSvc PermissionService) ItemsService {
	return &itemsService{
		repo:    repo,
		permSvc: permSvc,
	}
}

func (s *itemsService) Create(ctx context.Context, usr *model.User, req dto.CreateItemRequest) (*model.Item, error) {
	// cek permission role
	if !s.permSvc.CanCreateMasterData(usr.Role) {
		return nil, errors.New("forbidden: cannot create category")
	}

	createdBy := usr.ID

	// mapping
	items := &model.Item{
		CategoryID:   req.CategoryID,
		RackID:       req.RackID,
		SKU:          req.SKU,
		Name:         req.Name,
		Description:  req.Description,
		Unit:         req.Unit,
		Price:        req.Price,
		Cost:         req.Cost,
		Stock:        req.Stock,
		MinimumStock: req.MinimumStock,
		Weight:       req.Weight,
		Dimensions:   req.Dimensions,
		IsActive:     true,
		CreatedBy:    &createdBy,
	}

	return s.repo.Create(ctx, items)
}

func (s *itemsService) FindAll(page, limit int) (*[]model.Item, *dto.Pagination, error) {
	items, total, err := s.repo.Lists(page, limit)
	if err != nil {
		return nil, nil, err
	}

	pagination := dto.Pagination{
		Page:       page,
		Limit:      limit,
		TotalPages: utils.TotalPage(limit, int64(total)),
	}

	return &items, &pagination, nil
}

func (s *itemsService) Update(ctx context.Context, user *model.User, id int, req dto.UpdateItemRequest) (*model.Item, error) {

	if !s.permSvc.CanUpdateMasterData(user.Role) {
		return nil, errors.New("forbidden: cannot update item")
	}

	return s.repo.Update(ctx, id, req)
}

func (s *itemsService) FindByID(ctx context.Context, id int, usr *model.User) (*model.Item, error) {

	if !s.permSvc.CanReadMasterData(usr.Role) {
		return nil, errors.New("forbidden: cannot access items")
	}

	return s.repo.FindByID(ctx, id)
}

func (s *itemsService) Delete(ctx context.Context, usr *model.User, id int) error {
	if !s.permSvc.CanDeleteMasterData(usr.Role) {
		return errors.New("forbidden: cannont delete item")
	}

	return s.repo.Delete(ctx, id)
}
