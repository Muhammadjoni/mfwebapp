package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/muhammadjoni/mfwebapp/internal/domain/entity"
)

type SellerRepository interface {
	Create(ctx context.Context, seller *entity.Seller) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Seller, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Seller, error)
	Update(ctx context.Context, seller *entity.Seller) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.SellerStatus) error
	List(ctx context.Context, filter SellerFilter) ([]*entity.Seller, int64, error)
}

type SellerFilter struct {
	Status entity.SellerStatus
	Search string
	Page   int
	Limit  int
}
