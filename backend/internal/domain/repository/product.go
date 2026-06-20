package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/muhammadjoni/mfwebapp/internal/domain/entity"
)

type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Product, error)
	Update(ctx context.Context, product *entity.Product) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.ProductStatus) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter ProductFilter) ([]*entity.Product, int64, error)
	IncrementViewCount(ctx context.Context, id uuid.UUID) error
	DecrementStock(ctx context.Context, id uuid.UUID, qty int) error
	IncrementSoldCount(ctx context.Context, id uuid.UUID, qty int) error
}

type CategoryRepository interface {
	Create(ctx context.Context, cat *entity.Category) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Category, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Category, error)
	Update(ctx context.Context, cat *entity.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*entity.Category, error)
}

type ProductFilter struct {
	CategoryID *uuid.UUID
	SellerID   *uuid.UUID
	Status     entity.ProductStatus
	MinPrice   float64
	MaxPrice   float64
	Search     string
	Tags       []string
	SortBy     string
	SortDir    string
	Page       int
	Limit      int
}
