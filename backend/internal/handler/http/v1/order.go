package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/muhammadjoni/mfwebapp/internal/domain/entity"
	"github.com/muhammadjoni/mfwebapp/internal/handler/middleware"
	"github.com/muhammadjoni/mfwebapp/internal/service"
	"github.com/muhammadjoni/mfwebapp/pkg/response"
)

type OrderHandler struct {
	orderSvc *service.OrderService
}

func NewOrderHandler(orderSvc *service.OrderService) *OrderHandler {
	return &OrderHandler{orderSvc: orderSvc}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.ContextKeyUserID).(uuid.UUID)
	if !ok {
		response.Unauthorized(w)
		return
	}

	var req struct {
		Items []struct {
			ProductID string `json:"product_id"`
			VariantID string `json:"variant_id"`
			Quantity  int    `json:"quantity"`
		} `json:"items"`
		ShippingAddress entity.Address `json:"shipping_address"`
		Notes           string         `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	var items []service.OrderItemInput
	for _, i := range req.Items {
		pid, err := uuid.Parse(i.ProductID)
		if err != nil {
			response.BadRequest(w, "invalid product_id: "+i.ProductID)
			return
		}
		item := service.OrderItemInput{ProductID: pid, Quantity: i.Quantity}
		if i.VariantID != "" {
			if vid, err := uuid.Parse(i.VariantID); err == nil {
				item.VariantID = &vid
			}
		}
		items = append(items, item)
	}

	order, err := h.orderSvc.Create(r.Context(), service.CreateOrderInput{
		UserID:          userID,
		Items:           items,
		ShippingAddress: req.ShippingAddress,
		Notes:           req.Notes,
	})
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusCreated, order)
}

func (h *OrderHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	adminID, ok := r.Context().Value(middleware.ContextKeyUserID).(uuid.UUID)
	if !ok {
		response.Unauthorized(w)
		return
	}
	orderID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid order id")
		return
	}

	var req struct {
		Status string `json:"status"`
		Note   string `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if err := h.orderSvc.ChangeStatus(r.Context(), orderID, entity.OrderStatus(req.Status), adminID, req.Note, true); err != nil {
		switch err {
		case service.ErrOrderNotFound:
			response.NotFound(w, "order")
		case service.ErrInvalidStatusChange:
			response.BadRequest(w, "invalid status transition")
		default:
			response.InternalError(w)
		}
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "order status updated"})
}
