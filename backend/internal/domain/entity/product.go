package entity

import (
	"time"

	"github.com/google/uuid"
)

type ProductStatus string

const (
	ProductStatusDraft    ProductStatus = "draft"
	ProductStatusPending  ProductStatus = "pending"
	ProductStatusActive   ProductStatus = "active"
	ProductStatusInactive ProductStatus = "inactive"
	ProductStatusRejected ProductStatus = "rejected"
)

type Category struct {
	ID          uuid.UUID  `json:"id"`
	ParentID    *uuid.UUID `json:"parent_id"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Description string     `json:"description"`
	ImageURL    string     `json:"image_url"`
	SortOrder   int        `json:"sort_order"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Product struct {
	ID             uuid.UUID         `json:"id"`
	SellerID       uuid.UUID         `json:"seller_id"`
	CategoryID     uuid.UUID         `json:"category_id"`
	Name           string            `json:"name"`
	Slug           string            `json:"slug"`
	Description    string            `json:"description"`
	ShortDesc      string            `json:"short_description"`
	BasePrice      float64           `json:"base_price"`
	SalePrice      *float64          `json:"sale_price"`
	Currency       string            `json:"currency"`
	SKU            string            `json:"sku"`
	Stock          int               `json:"stock"`
	Status         ProductStatus     `json:"status"`
	Images         []string          `json:"images"`
	Tags           []string          `json:"tags"`
	Specifications map[string]string `json:"specifications"`
	Weight         float64           `json:"weight"`
	Dimensions     Dimensions        `json:"dimensions"`
	ViewCount      int               `json:"view_count"`
	SoldCount      int               `json:"sold_count"`
	Rating         float64           `json:"rating"`
	ReviewCount    int               `json:"review_count"`
	FeaturedAt     *time.Time        `json:"featured_at"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

type Dimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type ProductVariant struct {
	ID         uuid.UUID         `json:"id"`
	ProductID  uuid.UUID         `json:"product_id"`
	Name       string            `json:"name"`
	SKU        string            `json:"sku"`
	Price      float64           `json:"price"`
	Stock      int               `json:"stock"`
	Attributes map[string]string `json:"attributes"`
	ImageURL   string            `json:"image_url"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}
