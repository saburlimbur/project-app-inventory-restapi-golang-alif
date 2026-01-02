package service

import (
	"alfdwirhmn/inventory/dto"
	"alfdwirhmn/inventory/model"
	"alfdwirhmn/inventory/repository"
	"alfdwirhmn/inventory/utils"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*model.User, error)
	Login(ctx context.Context, req dto.LoginRequest, ip, ua string) (*dto.LoginResponse, error)
	Logout(ctx context.Context, token uuid.UUID) error
}

type authService struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
}

func NewAuthService(u repository.UserRepository, s repository.SessionRepository) AuthService {
	return &authService{u, s}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*model.User, error) {

	// hanya boleh staff & admin
	if req.Role != "staff" && req.Role != "admin" {
		return nil, errors.New("role not allowed")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hash),
		FullName:     req.FullName,
		Role:         req.Role,
		IsActive:     true,
	}

	return s.userRepo.Create(ctx, user)
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest, ip, ua string) (*dto.LoginResponse, error) {

	user, err := s.userRepo.FindByIdentifier(ctx, req.Identifier)
	if err != nil || !user.IsActive {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	token := uuid.New()
	expiredAt := time.Now().Add(24 * time.Hour)

	session := &model.Session{
		UserID:    user.ID,
		Token:     token,
		ExpiredAt: expiredAt,
		IPAddress: ip,
		UserAgent: ua,
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token:     token.String(),
		ExpiredAt: expiredAt.Format(time.RFC3339), // string
		User:      dto.ToUserResponseDTO(user),
	}, nil
}

func (s *authService) Logout(ctx context.Context, token uuid.UUID) error {
	return s.sessionRepo.Revoke(ctx, token)
}
