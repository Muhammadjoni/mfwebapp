package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/muhammadjoni/mfwebapp/internal/domain/entity"
	"github.com/muhammadjoni/mfwebapp/internal/domain/repository"
)

var (
	ErrOrderNotFound       = errors.New("order not found")
	ErrInvalidStatusChange = errors.New("invalid order status transition")
)

// allowedTransitions defines valid state machine transitions for non-admin users.
var allowedTransitions = map[entity.OrderStatus][]entity.OrderStatus{
	entity.OrderStatusCreated:    {entity.OrderStatusPaid, entity.OrderStatusCancelled},
	entity.OrderStatusPaid:       {entity.OrderStatusProcessing, entity.OrderStatusRefunded},
	entity.OrderStatusProcessing: {entity.OrderStatusShipped, entity.OrderStatusCancelled},
	entity.OrderStatusShipped:    {entity.OrderStatusDelivered},
	entity.OrderStatusDelivered:  {entity.OrderStatusClosed, entity.OrderStatusRefunded},
}

type OrderService struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
	cartRepo    repository.CartRepository
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
	cartRepo repository.CartRepository,
) *OrderService {
	return &OrderService{orderRepo: orderRepo, productRepo: productRepo, cartRepo: cartRepo}
}

type OrderItemInput struct {
	ProductID uuid.UUID
	VariantID *uuid.UUID
	Quantity  int
}

type CreateOrderInput struct {
	UserID          uuid.UUID
	Items           []OrderItemInput
	ShippingAddress entity.Address
	Notes           string
}

func (s *OrderService) Create(ctx context.Context, in CreateOrderInput) (*entity.Order, error) {
	var items []entity.OrderItem
	var subtotal float64

	for _, i := range in.Items {
		product, err := s.productRepo.GetByID(ctx, i.ProductID)
		if err != nil {
			return nil, ErrProductNotFound
		}

		price := product.BasePrice
		if product.SalePrice != nil {
			price = *product.SalePrice
		}

		lineTotal := price * float64(i.Quantity)
		subtotal += lineTotal

		items = append(items, entity.OrderItem{
			ID:         uuid.New(),
			ProductID:  i.ProductID,
			VariantID:  i.VariantID,
			SellerID:   product.SellerID,
			Name:       product.Name,
			SKU:        product.SKU,
			ImageURL:   firstImage(product.Images),
			Quantity:   i.Quantity,
			UnitPrice:  price,
			TotalPrice: lineTotal,
		})
	}

	order := &entity.Order{
		ID:              uuid.New(),
		UserID:          in.UserID,
		Status:          entity.OrderStatusCreated,
		Items:           items,
		SubTotal:        subtotal,
		ShippingCost:    0,
		Tax:             0,
		Total:           subtotal,
		Currency:        "USD",
		ShippingAddress: in.ShippingAddress,
		Notes:           in.Notes,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) ChangeStatus(ctx context.Context, orderID uuid.UUID, newStatus entity.OrderStatus, changedBy uuid.UUID, note string, adminOverride bool) error {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return ErrOrderNotFound
	}

	if !adminOverride && !isTransitionAllowed(order.Status, newStatus) {
		return ErrInvalidStatusChange
	}

	return s.orderRepo.UpdateStatus(ctx, orderID, newStatus, changedBy, note)
}

func isTransitionAllowed(from, to entity.OrderStatus) bool {
	for _, s := range allowedTransitions[from] {
		if s == to {
			return true
		}
	}
	return false
}

func (s *OrderService) GetByID(ctx context.Context, orderID, userID uuid.UUID) (*entity.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, ErrOrderNotFound
	}
	if order.UserID != userID {
		return nil, ErrOrderNotFound
	}
	return order, nil
}

func (s *OrderService) ListByUser(ctx context.Context, userID uuid.UUID, page, limit int) ([]*entity.Order, int64, error) {
	return s.orderRepo.GetByUserID(ctx, userID, page, limit)
}

func (s *OrderService) ListAll(ctx context.Context, filter repository.OrderFilter) ([]*entity.Order, int64, error) {
	return s.orderRepo.List(ctx, filter)
}

func firstImage(images []string) string {
	if len(images) > 0 {
		return images[0]
	}
	return ""
}
