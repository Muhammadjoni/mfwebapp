package entity

import (
	"time"

	"github.com/google/uuid"
)

type PaymentProvider string

const (
	PaymentProviderStripe PaymentProvider = "stripe"
	PaymentProviderVisa   PaymentProvider = "visa"
	PaymentProviderAlif   PaymentProvider = "alif"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSucceeded PaymentStatus = "succeeded"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

type Payment struct {
	ID               uuid.UUID         `json:"id"`
	OrderID          uuid.UUID         `json:"order_id"`
	UserID           uuid.UUID         `json:"user_id"`
	Provider         PaymentProvider   `json:"provider"`
	Status           PaymentStatus     `json:"status"`
	Amount           float64           `json:"amount"`
	Currency         string            `json:"currency"`
	ExternalID       string            `json:"external_id"`
	IdempotencyKey   string            `json:"idempotency_key"`
	ProviderMetadata map[string]string `json:"provider_metadata"`
	FailureReason    string            `json:"failure_reason"`
	RefundedAmount   float64           `json:"refunded_amount"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}

type PaymentIntent struct {
	ClientSecret string          `json:"client_secret"`
	PaymentID    uuid.UUID       `json:"payment_id"`
	Amount       float64         `json:"amount"`
	Currency     string          `json:"currency"`
	Provider     PaymentProvider `json:"provider"`
}
