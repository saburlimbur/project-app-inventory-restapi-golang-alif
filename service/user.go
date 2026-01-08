package service

import (
	"alfdwirhmn/inventory/dto"
	"alfdwirhmn/inventory/model"
	"alfdwirhmn/inventory/repository"
	"alfdwirhmn/inventory/utils"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Create(ctx context.Context, currentUser *model.User, req dto.CreateUserRequest) (*model.User, error)
	FindAll(page, limit int) (*[]model.User, *dto.Pagination, error)
	Update(ctx context.Context, currentUser *model.User, id int, req dto.UpdateUserRequest) error
	Delete(ctx context.Context, currentUser *model.User, id int) error
	Detail(ctx context.Context, id int) (*model.User, error)
}
type userService struct {
	repo    repository.UserRepository
	permSvc PermissionService
}

func NewUserService(repo repository.UserRepository, permSvc PermissionService) UserService {
	return &userService{
		repo:    repo,
		permSvc: permSvc,
	}
}

func (s *userService) Create(ctx context.Context, currentUser *model.User, req dto.CreateUserRequest) (*model.User, error) {
	// cek permission
	if err := s.permSvc.CanCreateUser(currentUser.Role, req.Role); err != nil {
		return nil, err
	}

	// cek email
	exists, err := s.repo.IsEmailExists(ctx, req.Email, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// cek username
	exists, err = s.repo.IsUsernameExists(ctx, req.Username, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		Role:         req.Role,
		IsActive:     true,
	}

	return s.repo.Create(ctx, user)
}

func (s *userService) FindAll(page, limit int) (*[]model.User, *dto.Pagination, error) {
	users, total, err := s.repo.Lists(page, limit)

	if err != nil {
		return nil, nil, err
	}

	pagination := dto.Pagination{
		Page:       page,
		Limit:      limit,
		TotalPages: utils.TotalPage(limit, int64(total)),
		TotalRows:  total,
	}

	return &users, &pagination, nil
}

func (s *userService) Detail(ctx context.Context, id int) (*model.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *userService) Update(ctx context.Context, currentUser *model.User, id int, req dto.UpdateUserRequest) error {

	if currentUser.Role != "admin" && currentUser.Role != "super_admin" {
		return errors.New("forbidden")
	}

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if req.Username != nil {
		user.Username = *req.Username
	}

	if req.Email != nil {
		user.Email = *req.Email
	}

	if req.FullName != nil {
		user.FullName = *req.FullName
	}

	if req.Role != nil {
		user.Role = *req.Role
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	return s.repo.Update(ctx, user)
}

func (s *userService) Delete(ctx context.Context, currentUser *model.User, id int) error {
	if currentUser.Role != "admin" && currentUser.Role != "super_admin" {
		return errors.New("forbidden")
	}

	return s.repo.Delete(ctx, id)
}
