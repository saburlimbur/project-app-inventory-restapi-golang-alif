package service

import (
	"alfdwirhmn/inventory/dto"
	"alfdwirhmn/inventory/model"
	"alfdwirhmn/inventory/repository"
	"alfdwirhmn/inventory/utils"
	"context"
	"errors"
)

type WarehouseService interface {
	Create(ctx context.Context, usr *model.User, req dto.CreateWarehouseRequest) (*model.Warehouse, error)
	FindAll(page, limit int) (*[]model.Warehouse, *dto.Pagination, error)
	Update(ctx context.Context, usr *model.User, id int, req dto.UpdateWarehouseRequest) (*model.Warehouse, error)
	Delete(ctx context.Context, usr *model.User, id int) error
}

type warehouseService struct {
	repo    repository.WarehouseRepository
	permSvc PermissionService
}

func NewWarehouseService(repo repository.WarehouseRepository, permSvc PermissionService) WarehouseService {
	return &warehouseService{
		repo:    repo,
		permSvc: permSvc,
	}
}

func (w *warehouseService) Create(ctx context.Context, usr *model.User, req dto.CreateWarehouseRequest) (*model.Warehouse, error) {
	if !w.permSvc.CanCreateMasterData(usr.Role) {
		return nil, errors.New("forbidden: cannon create warehouse")
	}

	createdBy := usr.ID

	warehouse := &model.Warehouse{
		Code:       req.Code,
		Name:       req.Name,
		Address:    req.Address,
		City:       req.City,
		Province:   req.Province,
		PostalCode: req.PostalCode,
		Phone:      req.Phone,
		IsActive:   true,
		CreatedBy:  &createdBy,
	}

	return w.repo.Create(ctx, warehouse)
}

func (w *warehouseService) FindAll(page, limit int) (*[]model.Warehouse, *dto.Pagination, error) {
	warehouses, total, err := w.repo.Lists(page, limit)
	if err != nil {
		return nil, nil, err
	}

	pagination := dto.Pagination{
		Page:       page,
		Limit:      limit,
		TotalPages: utils.TotalPage(limit, int64(total)),
		TotalRows:  total,
	}

	return &warehouses, &pagination, nil
}

func (c *warehouseService) Update(ctx context.Context, usr *model.User, id int, req dto.UpdateWarehouseRequest) (*model.Warehouse, error) {

	// permission
	if !c.permSvc.CanUpdateMasterData(usr.Role) {
		return nil, errors.New("forbidden: cannot update category")
	}

	payload := &model.Warehouse{
		Code:       req.Code,
		Name:       req.Name,
		Address:    req.Address,
		City:       req.City,
		Province:   req.Province,
		PostalCode: req.PostalCode,
		Phone:      req.Phone,
		IsActive:   req.IsActive,
	}

	return c.repo.Update(ctx, id, payload)
}

func (c *warehouseService) Delete(ctx context.Context, usr *model.User, id int) error {
	if !c.permSvc.CanDeleteMasterData(usr.Role) {
		return errors.New("forbidden: cannot delete category")
	}

	return c.repo.Delete(ctx, id)
}
