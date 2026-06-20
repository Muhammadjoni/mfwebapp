package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/muhammadjoni/mfwebapp/internal/domain/entity"
	"github.com/muhammadjoni/mfwebapp/internal/domain/repository"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

type ProductService struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
}

func NewProductService(productRepo repository.ProductRepository, categoryRepo repository.CategoryRepository) *ProductService {
	return &ProductService{productRepo: productRepo, categoryRepo: categoryRepo}
}

type CreateProductInput struct {
	SellerID       uuid.UUID
	CategoryID     uuid.UUID
	Name           string
	Description    string
	ShortDesc      string
	BasePrice      float64
	Currency       string
	SKU            string
	Stock          int
	Images         []string
	Tags           []string
	Specifications map[string]string
	Weight         float64
}

func (s *ProductService) Create(ctx context.Context, in CreateProductInput) (*entity.Product, error) {
	product := &entity.Product{
		ID:             uuid.New(),
		SellerID:       in.SellerID,
		CategoryID:     in.CategoryID,
		Name:           in.Name,
		Slug:           slugify(in.Name),
		Description:    in.Description,
		ShortDesc:      in.ShortDesc,
		BasePrice:      in.BasePrice,
		Currency:       in.Currency,
		SKU:            in.SKU,
		Stock:          in.Stock,
		Status:         entity.ProductStatusPending,
		Images:         in.Images,
		Tags:           in.Tags,
		Specifications: in.Specifications,
		Weight:         in.Weight,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrProductNotFound
	}
	go s.productRepo.IncrementViewCount(context.Background(), id) //nolint:errcheck
	return product, nil
}

func (s *ProductService) List(ctx context.Context, filter repository.ProductFilter) ([]*entity.Product, int64, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	return s.productRepo.List(ctx, filter)
}

func (s *ProductService) Approve(ctx context.Context, id uuid.UUID) error {
	return s.productRepo.UpdateStatus(ctx, id, entity.ProductStatusActive)
}

func (s *ProductService) Reject(ctx context.Context, id uuid.UUID) error {
	return s.productRepo.UpdateStatus(ctx, id, entity.ProductStatusRejected)
}

func slugify(s string) string {
	slug := strings.ToLower(strings.ReplaceAll(s, " ", "-"))
	return slug + "-" + uuid.New().String()[:8]
}
