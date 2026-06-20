package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/muhammadjoni/mfwebapp/internal/domain/entity"
	"github.com/muhammadjoni/mfwebapp/internal/domain/repository"
)

var (
	ErrPaymentNotFound = errors.New("payment not found")
	ErrUnknownProvider = errors.New("unsupported payment provider")
)

// PaymentProvider is the abstraction all payment integrations must satisfy.
type PaymentProvider interface {
	Charge(ctx context.Context, req ChargeRequest) (*ChargeResponse, error)
	Refund(ctx context.Context, externalID string, amount float64) error
	VerifyWebhook(payload []byte, signature string) (WebhookEvent, error)
}

type ChargeRequest struct {
	Amount         float64
	Currency       string
	IdempotencyKey string
	OrderID        string
	CustomerEmail  string
}

type ChargeResponse struct {
	ExternalID   string
	ClientSecret string
	Status       entity.PaymentStatus
}

type WebhookEvent struct {
	ExternalID string
	Status     entity.PaymentStatus
	Provider   entity.PaymentProvider
}

type PaymentService struct {
	paymentRepo repository.PaymentRepository
	orderRepo   repository.OrderRepository
	providers   map[entity.PaymentProvider]PaymentProvider
}

func NewPaymentService(
	paymentRepo repository.PaymentRepository,
	orderRepo repository.OrderRepository,
	providers map[entity.PaymentProvider]PaymentProvider,
) *PaymentService {
	return &PaymentService{paymentRepo: paymentRepo, orderRepo: orderRepo, providers: providers}
}

func (s *PaymentService) Initiate(ctx context.Context, orderID, userID uuid.UUID, provider entity.PaymentProvider, email string) (*entity.PaymentIntent, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	p, ok := s.providers[provider]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownProvider, provider)
	}

	idempotencyKey := fmt.Sprintf("%s:%s", orderID, provider)

	resp, err := p.Charge(ctx, ChargeRequest{
		Amount:         order.Total,
		Currency:       order.Currency,
		IdempotencyKey: idempotencyKey,
		OrderID:        orderID.String(),
		CustomerEmail:  email,
	})
	if err != nil {
		return nil, err
	}

	payment := &entity.Payment{
		ID:             uuid.New(),
		OrderID:        orderID,
		UserID:         userID,
		Provider:       provider,
		Status:         entity.PaymentStatusPending,
		Amount:         order.Total,
		Currency:       order.Currency,
		ExternalID:     resp.ExternalID,
		IdempotencyKey: idempotencyKey,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, err
	}

	return &entity.PaymentIntent{
		ClientSecret: resp.ClientSecret,
		PaymentID:    payment.ID,
		Amount:       order.Total,
		Currency:     order.Currency,
		Provider:     provider,
	}, nil
}

func (s *PaymentService) HandleWebhook(ctx context.Context, provider entity.PaymentProvider, payload []byte, signature string) error {
	p, ok := s.providers[provider]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownProvider, provider)
	}

	event, err := p.VerifyWebhook(payload, signature)
	if err != nil {
		return fmt.Errorf("webhook verification failed: %w", err)
	}

	payment, err := s.paymentRepo.GetByExternalID(ctx, provider, event.ExternalID)
	if err != nil {
		return ErrPaymentNotFound
	}

	payment.Status = event.Status
	payment.UpdatedAt = time.Now()

	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return err
	}

	if event.Status == entity.PaymentStatusSucceeded {
		_ = s.orderRepo.UpdateStatus(ctx, payment.OrderID, entity.OrderStatusPaid, uuid.Nil, "payment confirmed via webhook")
	}

	return nil
}
