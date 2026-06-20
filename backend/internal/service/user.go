package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/muhammadjoni/mfwebapp/internal/domain/entity"
	"github.com/muhammadjoni/mfwebapp/internal/domain/repository"
)

var ErrUserNotFound = errors.New("user not found")

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) List(ctx context.Context, filter repository.UserFilter) ([]*entity.User, int64, error) {
	return s.userRepo.List(ctx, filter)
}

func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.UserStatus) error {
	return s.userRepo.UpdateStatus(ctx, id, status)
}

func (s *UserService) UpdateRole(ctx context.Context, id uuid.UUID, role entity.Role) error {
	return s.userRepo.UpdateRole(ctx, id, role)
}
