package service

import (
	"alfdwirhmn/inventory/dto"
	"alfdwirhmn/inventory/model"
	"alfdwirhmn/inventory/repository"
	"alfdwirhmn/inventory/utils"
	"context"
	"errors"
)

type RacksService interface {
	Create(ctx context.Context, usr *model.User, req dto.CreateRackRequest) (*model.Racks, error)
	FindAll(page, limit int) ([]model.Racks, *dto.Pagination, error)
	FindById(id int, usr *model.User) (*model.Racks, error)
	Update(ctx context.Context, usr *model.User, id int, req dto.UpdateRackRequest) (*model.Racks, error)
	Delete(ctx context.Context, usr *model.User, id int) error
}

type racksService struct {
	repo    repository.RacksRepository
	permSvc PermissionService
}

func NewRacksService(repo repository.RacksRepository, permSvc PermissionService) RacksService {
	return &racksService{
		repo:    repo,
		permSvc: permSvc,
	}
}

func (s *racksService) Create(ctx context.Context, usr *model.User, req dto.CreateRackRequest) (*model.Racks, error) {
	if !s.permSvc.CanCreateMasterData(usr.Role) {
		return nil, errors.New("forbidden: cannot create rack")
	}

	createdBy := usr.ID

	rack := &model.Racks{
		WarehouseID: req.WarehouseID,
		Code:        req.Code,
		Name:        req.Name,
		Location:    req.Location,
		Capacity:    req.Capacity,
		Description: req.Description,
		IsActive:    true,
		CreatedBy:   &createdBy,
	}

	return s.repo.Create(ctx, rack)
}

func (s *racksService) FindAll(page, limit int) ([]model.Racks, *dto.Pagination, error) {
	racks, total, err := s.repo.Lists(page, limit)
	if err != nil {
		return nil, nil, err
	}

	pagination := &dto.Pagination{
		Page:       page,
		Limit:      limit,
		TotalPages: utils.TotalPage(limit, int64(total)),
		TotalRows:  total,
	}

	return racks, pagination, nil
}

func (s *racksService) FindById(id int, usr *model.User) (*model.Racks, error) {
	if !s.permSvc.CanReadMasterData(usr.Role) {
		return nil, errors.New("forbidden: cannon access racks")
	}

	return s.repo.DetailById(id)
}

func (s *racksService) Update(ctx context.Context, usr *model.User, id int, req dto.UpdateRackRequest) (*model.Racks, error) {
	if !s.permSvc.CanUpdateMasterData(usr.Role) {
		return nil, errors.New("forbidden: cannot update rack")
	}

	payload := &model.Racks{
		Code:        req.Code,
		Name:        req.Name,
		Location:    req.Location,
		Capacity:    req.Capacity,
		Description: req.Description,
		IsActive:    req.IsActive,
	}

	return s.repo.Update(ctx, id, payload)
}

func (s *racksService) Delete(ctx context.Context, usr *model.User, id int) error {
	if !s.permSvc.CanDeleteMasterData(usr.Role) {
		return errors.New("forbidden: cannot delete rack")
	}

	return s.repo.Delete(ctx, id)
}
