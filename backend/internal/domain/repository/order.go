package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/muhammadjoni/mfwebapp/internal/domain/entity"
)

type OrderRepository interface {
	Create(ctx context.Context, order *entity.Order) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Order, error)
	Update(ctx context.Context, order *entity.Order) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.OrderStatus, changedBy uuid.UUID, note string) error
	List(ctx context.Context, filter OrderFilter) ([]*entity.Order, int64, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]*entity.Order, int64, error)
}

type PaymentRepository interface {
	Create(ctx context.Context, payment *entity.Payment) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error)
	GetByOrderID(ctx context.Context, orderID uuid.UUID) (*entity.Payment, error)
	GetByExternalID(ctx context.Context, provider entity.PaymentProvider, externalID string) (*entity.Payment, error)
	Update(ctx context.Context, payment *entity.Payment) error
	List(ctx context.Context, filter PaymentFilter) ([]*entity.Payment, int64, error)
}

type CartRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Cart, error)
	UpsertItem(ctx context.Context, cartID, productID uuid.UUID, variantID *uuid.UUID, qty int) error
	RemoveItem(ctx context.Context, cartID, productID uuid.UUID, variantID *uuid.UUID) error
	Clear(ctx context.Context, cartID uuid.UUID) error
	GetOrCreate(ctx context.Context, userID uuid.UUID) (*entity.Cart, error)
}

type OrderFilter struct {
	UserID   *uuid.UUID
	SellerID *uuid.UUID
	Status   entity.OrderStatus
	Page     int
	Limit    int
}

type PaymentFilter struct {
	UserID   *uuid.UUID
	Provider entity.PaymentProvider
	Status   entity.PaymentStatus
	Page     int
	Limit    int
}
