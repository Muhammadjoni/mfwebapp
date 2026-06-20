package entity

import (
	"time"

	"github.com/google/uuid"
)

type SellerStatus string

const (
	SellerStatusPending   SellerStatus = "pending"
	SellerStatusApproved  SellerStatus = "approved"
	SellerStatusRejected  SellerStatus = "rejected"
	SellerStatusSuspended SellerStatus = "suspended"
)

type Seller struct {
	ID             uuid.UUID    `json:"id"`
	UserID         uuid.UUID    `json:"user_id"`
	BusinessName   string       `json:"business_name"`
	BusinessEmail  string       `json:"business_email"`
	BusinessPhone  string       `json:"business_phone"`
	Country        string       `json:"country"`
	Status         SellerStatus `json:"status"`
	CommissionRate float64      `json:"commission_rate"`
	Description    string       `json:"description"`
	LogoURL        string       `json:"logo_url"`
	VerifiedAt     *time.Time   `json:"verified_at"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

type Cart struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	Items     []CartItem `json:"items"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CartItem struct {
	ID        uuid.UUID  `json:"id"`
	CartID    uuid.UUID  `json:"cart_id"`
	ProductID uuid.UUID  `json:"product_id"`
	VariantID *uuid.UUID `json:"variant_id"`
	Quantity  int        `json:"quantity"`
	AddedAt   time.Time  `json:"added_at"`
}
