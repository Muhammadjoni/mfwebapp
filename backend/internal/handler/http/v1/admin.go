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

type AdminHandler struct {
	userSvc   *service.UserService
	sellerSvc *service.SellerService
	orderSvc  *service.OrderService
	productSvc *service.ProductService
}

func NewAdminHandler(
	userSvc *service.UserService,
	sellerSvc *service.SellerService,
	orderSvc *service.OrderService,
	productSvc *service.ProductService,
) *AdminHandler {
	return &AdminHandler{
		userSvc:   userSvc,
		sellerSvc: sellerSvc,
		orderSvc:  orderSvc,
		productSvc: productSvc,
	}
}

// Users

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	filter := repository.UserFilter{
		Search: q.Get("search"),
		Page:   page,
		Limit:  limit,
	}
	if s := q.Get("status"); s != "" {
		filter.Status = entity.UserStatus(s)
	}
	if role := q.Get("role"); role != "" {
		filter.Role = entity.Role(role)
	}

	users, total, err := h.userSvc.List(r.Context(), filter)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSONWithMeta(w, http.StatusOK, users, &response.Meta{
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: (total + int64(limit) - 1) / int64(limit),
	})
}

func (h *AdminHandler) UpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid user id")
		return
	}
	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if err := h.userSvc.UpdateStatus(r.Context(), id, entity.UserStatus(req.Status)); err != nil {
		response.NotFound(w, "user")
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "user status updated"})
}

func (h *AdminHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid user id")
		return
	}
	var req struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if err := h.userSvc.UpdateRole(r.Context(), id, entity.Role(req.Role)); err != nil {
		response.NotFound(w, "user")
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "user role updated"})
}

// Sellers

func (h *AdminHandler) ListSellers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	filter := repository.SellerFilter{
		Search: q.Get("search"),
		Page:   page,
		Limit:  limit,
	}
	if s := q.Get("status"); s != "" {
		filter.Status = entity.SellerStatus(s)
	}

	sellers, total, err := h.sellerSvc.List(r.Context(), filter)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSONWithMeta(w, http.StatusOK, sellers, &response.Meta{
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: (total + int64(limit) - 1) / int64(limit),
	})
}

func (h *AdminHandler) ApproveSeller(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid seller id")
		return
	}
	if err := h.sellerSvc.Approve(r.Context(), id); err != nil {
		response.NotFound(w, "seller")
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "seller approved"})
}

func (h *AdminHandler) RejectSeller(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid seller id")
		return
	}
	if err := h.sellerSvc.Reject(r.Context(), id); err != nil {
		response.NotFound(w, "seller")
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "seller rejected"})
}

// Orders (admin view — all orders)

func (h *AdminHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	filter := repository.OrderFilter{
		Page:  page,
		Limit: limit,
	}
	if s := q.Get("status"); s != "" {
		filter.Status = entity.OrderStatus(s)
	}

	orders, total, err := h.orderSvc.ListAll(r.Context(), filter)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSONWithMeta(w, http.StatusOK, orders, &response.Meta{
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: (total + int64(limit) - 1) / int64(limit),
	})
}

// Payments

func (h *AdminHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	filter := repository.ProductFilter{
		Search: q.Get("search"),
		Page:   page,
		Limit:  limit,
	}
	if s := q.Get("status"); s != "" {
		filter.Status = entity.ProductStatus(s)
	}

	products, total, err := h.productSvc.List(r.Context(), filter)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSONWithMeta(w, http.StatusOK, products, &response.Meta{
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: (total + int64(limit) - 1) / int64(limit),
	})
}

// Me — get current user profile

func (h *AdminHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.ContextKeyUserID).(uuid.UUID)
	if !ok {
		response.Unauthorized(w)
		return
	}
	user, err := h.userSvc.GetByID(r.Context(), userID)
	if err != nil {
		response.NotFound(w, "user")
		return
	}
	response.JSON(w, http.StatusOK, user)
}
