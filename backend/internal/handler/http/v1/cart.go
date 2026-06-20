package v1

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/muhammadjoni/mfwebapp/internal/handler/middleware"
	"github.com/muhammadjoni/mfwebapp/internal/service"
	"github.com/muhammadjoni/mfwebapp/pkg/response"
)

type CartHandler struct {
	cartSvc *service.CartService
}

func NewCartHandler(cartSvc *service.CartService) *CartHandler {
	return &CartHandler{cartSvc: cartSvc}
}

func (h *CartHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.ContextKeyUserID).(uuid.UUID)
	if !ok {
		response.Unauthorized(w)
		return
	}
	cart, err := h.cartSvc.GetCart(r.Context(), userID)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, cart)
}

func (h *CartHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.ContextKeyUserID).(uuid.UUID)
	if !ok {
		response.Unauthorized(w)
		return
	}

	var req struct {
		ProductID string `json:"product_id"`
		VariantID string `json:"variant_id"`
		Quantity  int    `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if req.Quantity < 1 {
		req.Quantity = 1
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		response.BadRequest(w, "invalid product_id")
		return
	}

	var variantID *uuid.UUID
	if req.VariantID != "" {
		vid, err := uuid.Parse(req.VariantID)
		if err != nil {
			response.BadRequest(w, "invalid variant_id")
			return
		}
		variantID = &vid
	}

	cart, err := h.cartSvc.AddItem(r.Context(), userID, productID, variantID, req.Quantity)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, cart)
}

func (h *CartHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.ContextKeyUserID).(uuid.UUID)
	if !ok {
		response.Unauthorized(w)
		return
	}

	var req struct {
		ProductID string `json:"product_id"`
		VariantID string `json:"variant_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		response.BadRequest(w, "invalid product_id")
		return
	}

	var variantID *uuid.UUID
	if req.VariantID != "" {
		vid, err := uuid.Parse(req.VariantID)
		if err == nil {
			variantID = &vid
		}
	}

	cart, err := h.cartSvc.RemoveItem(r.Context(), userID, productID, variantID)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, cart)
}

func (h *CartHandler) Clear(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.ContextKeyUserID).(uuid.UUID)
	if !ok {
		response.Unauthorized(w)
		return
	}
	if err := h.cartSvc.Clear(r.Context(), userID); err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "cart cleared"})
}
