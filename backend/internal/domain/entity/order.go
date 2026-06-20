package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusCreated    OrderStatus = "created"
	OrderStatusPaid       OrderStatus = "paid"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusClosed     OrderStatus = "closed"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRefunded   OrderStatus = "refunded"
)

type Order struct {
	ID              uuid.UUID   `json:"id"`
	UserID          uuid.UUID   `json:"user_id"`
	Status          OrderStatus `json:"status"`
	Items           []OrderItem `json:"items"`
	SubTotal        float64     `json:"subtotal"`
	ShippingCost    float64     `json:"shipping_cost"`
	Tax             float64     `json:"tax"`
	Total           float64     `json:"total"`
	Currency        string      `json:"currency"`
	ShippingAddress Address     `json:"shipping_address"`
	Notes           string      `json:"notes"`
	TrackingNumber  string      `json:"tracking_number"`
	PaymentID       *uuid.UUID  `json:"payment_id"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID         uuid.UUID  `json:"id"`
	OrderID    uuid.UUID  `json:"order_id"`
	ProductID  uuid.UUID  `json:"product_id"`
	VariantID  *uuid.UUID `json:"variant_id"`
	SellerID   uuid.UUID  `json:"seller_id"`
	Name       string     `json:"name"`
	SKU        string     `json:"sku"`
	ImageURL   string     `json:"image_url"`
	Quantity   int        `json:"quantity"`
	UnitPrice  float64    `json:"unit_price"`
	TotalPrice float64    `json:"total_price"`
}

type Address struct {
	FullName   string `json:"full_name"`
	Phone      string `json:"phone"`
	Country    string `json:"country"`
	City       string `json:"city"`
	State      string `json:"state"`
	Street     string `json:"street"`
	PostalCode string `json:"postal_code"`
	Apartment  string `json:"apartment"`
}

type OrderStatusHistory struct {
	ID        uuid.UUID   `json:"id"`
	OrderID   uuid.UUID   `json:"order_id"`
	Status    OrderStatus `json:"status"`
	ChangedBy uuid.UUID   `json:"changed_by"`
	Note      string      `json:"note"`
	CreatedAt time.Time   `json:"created_at"`
}
