package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/muhammadjoni/mfwebapp/internal/domain/entity"
	"github.com/muhammadjoni/mfwebapp/internal/domain/repository"
)

type CartService struct {
	cartRepo repository.CartRepository
}

func NewCartService(cartRepo repository.CartRepository) *CartService {
	return &CartService{cartRepo: cartRepo}
}

func (s *CartService) GetCart(ctx context.Context, userID uuid.UUID) (*entity.Cart, error) {
	return s.cartRepo.GetOrCreate(ctx, userID)
}

func (s *CartService) AddItem(ctx context.Context, userID, productID uuid.UUID, variantID *uuid.UUID, qty int) (*entity.Cart, error) {
	cart, err := s.cartRepo.GetOrCreate(ctx, userID)
	if err != nil {
		return nil, err
	}
	if err := s.cartRepo.UpsertItem(ctx, cart.ID, productID, variantID, qty); err != nil {
		return nil, err
	}
	return s.cartRepo.GetOrCreate(ctx, userID)
}

func (s *CartService) RemoveItem(ctx context.Context, userID, productID uuid.UUID, variantID *uuid.UUID) (*entity.Cart, error) {
	cart, err := s.cartRepo.GetOrCreate(ctx, userID)
	if err != nil {
		return nil, err
	}
	if err := s.cartRepo.RemoveItem(ctx, cart.ID, productID, variantID); err != nil {
		return nil, err
	}
	return s.cartRepo.GetOrCreate(ctx, userID)
}

func (s *CartService) Clear(ctx context.Context, userID uuid.UUID) error {
	cart, err := s.cartRepo.GetOrCreate(ctx, userID)
	if err != nil {
		return err
	}
	return s.cartRepo.Clear(ctx, cart.ID)
}
