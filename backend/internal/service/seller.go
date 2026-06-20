package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/muhammadjoni/mfwebapp/internal/domain/entity"
	"github.com/muhammadjoni/mfwebapp/internal/domain/repository"
)

var ErrSellerNotFound = errors.New("seller not found")

type SellerService struct {
	sellerRepo repository.SellerRepository
}

func NewSellerService(sellerRepo repository.SellerRepository) *SellerService {
	return &SellerService{sellerRepo: sellerRepo}
}

func (s *SellerService) List(ctx context.Context, filter repository.SellerFilter) ([]*entity.Seller, int64, error) {
	return s.sellerRepo.List(ctx, filter)
}

func (s *SellerService) Approve(ctx context.Context, id uuid.UUID) error {
	seller, err := s.sellerRepo.GetByID(ctx, id)
	if err != nil {
		return ErrSellerNotFound
	}
	now := time.Now()
	seller.Status = entity.SellerStatusApproved
	seller.VerifiedAt = &now
	seller.UpdatedAt = now
	return s.sellerRepo.Update(ctx, seller)
}

func (s *SellerService) Reject(ctx context.Context, id uuid.UUID) error {
	if err := s.sellerRepo.UpdateStatus(ctx, id, entity.SellerStatusRejected); err != nil {
		return ErrSellerNotFound
	}
	return nil
}

func (s *SellerService) Register(ctx context.Context, userID uuid.UUID, businessName, businessEmail, businessPhone, country, description string) (*entity.Seller, error) {
	seller := &entity.Seller{
		ID:             uuid.New(),
		UserID:         userID,
		BusinessName:   businessName,
		BusinessEmail:  businessEmail,
		BusinessPhone:  businessPhone,
		Country:        country,
		Description:    description,
		Status:         entity.SellerStatusPending,
		CommissionRate: 0.10,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := s.sellerRepo.Create(ctx, seller); err != nil {
		return nil, err
	}
	return seller, nil
}
