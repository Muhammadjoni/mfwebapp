package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/muhammadjoni/mfwebapp/internal/domain/entity"
	"github.com/muhammadjoni/mfwebapp/internal/domain/repository"
	"github.com/muhammadjoni/mfwebapp/internal/handler/middleware"
	"github.com/muhammadjoni/mfwebapp/internal/service"
	"github.com/muhammadjoni/mfwebapp/pkg/response"
)

type ProductHandler struct {
	productSvc *service.ProductService
}

func NewProductHandler(productSvc *service.ProductService) *ProductHandler {
	return &ProductHandler{productSvc: productSvc}
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))
	minPrice, _ := strconv.ParseFloat(q.Get("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(q.Get("max_price"), 64)

	filter := repository.ProductFilter{
		Status:  entity.ProductStatusActive,
		Search:  q.Get("search"),
		SortBy:  q.Get("sort_by"),
		SortDir: q.Get("sort_dir"),
		Page:    page,
		Limit:   limit,
	}
	if minPrice > 0 {
		filter.MinPrice = minPrice
	}
	if maxPrice > 0 {
		filter.MaxPrice = maxPrice
	}
	if catID := q.Get("category_id"); catID != "" {
		if id, err := uuid.Parse(catID); err == nil {
			filter.CategoryID = &id
		}
	}

	products, total, err := h.productSvc.List(r.Context(), filter)
	if err != nil {
		response.InternalError(w)
		return
	}

	lim := limit
	if lim < 1 {
		lim = 20
	}
	pg := page
	if pg < 1 {
		pg = 1
	}

	response.JSONWithMeta(w, http.StatusOK, products, &response.Meta{
		Page:       pg,
		Limit:      lim,
		TotalItems: total,
		TotalPages: (total + int64(lim) - 1) / int64(lim),
	})
}

func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid product id")
		return
	}
	product, err := h.productSvc.GetByID(r.Context(), id)
	if err != nil {
		response.NotFound(w, "product")
		return
	}
	response.JSON(w, http.StatusOK, product)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	sellerID, ok := r.Context().Value(middleware.ContextKeyUserID).(uuid.UUID)
	if !ok {
		response.Unauthorized(w)
		return
	}

	var req struct {
		CategoryID     string            `json:"category_id"`
		Name           string            `json:"name"`
		Description    string            `json:"description"`
		ShortDesc      string            `json:"short_description"`
		BasePrice      float64           `json:"base_price"`
		Currency       string            `json:"currency"`
		SKU            string            `json:"sku"`
		Stock          int               `json:"stock"`
		Images         []string          `json:"images"`
		Tags           []string          `json:"tags"`
		Specifications map[string]string `json:"specifications"`
		Weight         float64           `json:"weight"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	catID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		response.BadRequest(w, "invalid category_id")
		return
	}

	product, err := h.productSvc.Create(r.Context(), service.CreateProductInput{
		SellerID:       sellerID,
		CategoryID:     catID,
		Name:           req.Name,
		Description:    req.Description,
		ShortDesc:      req.ShortDesc,
		BasePrice:      req.BasePrice,
		Currency:       req.Currency,
		SKU:            req.SKU,
		Stock:          req.Stock,
		Images:         req.Images,
		Tags:           req.Tags,
		Specifications: req.Specifications,
		Weight:         req.Weight,
	})
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusCreated, product)
}

func (h *ProductHandler) Approve(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid product id")
		return
	}
	if err := h.productSvc.Approve(r.Context(), id); err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "product approved"})
}

func (h *ProductHandler) Reject(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid product id")
		return
	}
	if err := h.productSvc.Reject(r.Context(), id); err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "product rejected"})
}
