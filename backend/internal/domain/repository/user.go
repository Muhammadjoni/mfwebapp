package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/muhammadjoni/mfwebapp/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.UserStatus) error
	UpdateRole(ctx context.Context, id uuid.UUID, role entity.Role) error
	List(ctx context.Context, filter UserFilter) ([]*entity.User, int64, error)
	Delete(ctx context.Context, id uuid.UUID) error

	SaveRefreshToken(ctx context.Context, token *entity.RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*entity.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteUserRefreshTokens(ctx context.Context, userID uuid.UUID) error
}

type UserFilter struct {
	Role   entity.Role
	Status entity.UserStatus
	Search string
	Page   int
	Limit  int
}
