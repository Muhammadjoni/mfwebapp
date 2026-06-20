package v1

import (
	"encoding/json"
	"net/http"

	"github.com/muhammadjoni/mfwebapp/internal/domain/entity"
	"github.com/muhammadjoni/mfwebapp/internal/service"
	"github.com/muhammadjoni/mfwebapp/pkg/response"
)

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email     string          `json:"email"`
		Password  string          `json:"password"`
		FirstName string          `json:"first_name"`
		LastName  string          `json:"last_name"`
		Phone     string          `json:"phone"`
		Language  entity.Language `json:"language"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if req.Email == "" || req.Password == "" || req.FirstName == "" {
		response.BadRequest(w, "email, password, and first_name are required")
		return
	}
	if len(req.Password) < 8 {
		response.BadRequest(w, "password must be at least 8 characters")
		return
	}

	user, tokens, err := h.authSvc.Register(r.Context(), service.RegisterInput{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Language:  req.Language,
	})
	if err != nil {
		switch err {
		case service.ErrEmailTaken:
			response.Error(w, http.StatusConflict, "EMAIL_TAKEN", "email already in use")
		default:
			response.InternalError(w)
		}
		return
	}

	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"user":   user,
		"tokens": tokens,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	user, tokens, err := h.authSvc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			response.Error(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid email or password")
		case service.ErrUserBanned:
			response.Error(w, http.StatusForbidden, "ACCOUNT_BANNED", "your account has been suspended")
		default:
			response.InternalError(w)
		}
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"user":   user,
		"tokens": tokens,
	})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	tokens, err := h.authSvc.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "INVALID_TOKEN", "token expired or invalid")
		return
	}
	response.JSON(w, http.StatusOK, tokens)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	_ = h.authSvc.Logout(r.Context(), req.RefreshToken)
	response.JSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}
