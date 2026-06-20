package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/muhammadjoni/mfwebapp/internal/handler/middleware"
)

type Router struct {
	auth    *AuthHandler
	product *ProductHandler
	order   *OrderHandler
	authMW  *middleware.AuthMiddleware
	origins []string
	rateRPM int
}

func NewRouter(
	auth *AuthHandler,
	product *ProductHandler,
	order *OrderHandler,
	authMW *middleware.AuthMiddleware,
	origins []string,
	rateRPM int,
) *Router {
	return &Router{auth: auth, product: product, order: order, authMW: authMW, origins: origins, rateRPM: rateRPM}
}

func (rt *Router) Build() http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.StripSlashes)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   rt.origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(httprate.LimitByIP(rt.rateRPM, 60))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`)) //nolint:errcheck
	})

	r.Route("/api/v1", func(r chi.Router) {
		// Public
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", rt.auth.Register)
			r.Post("/login", rt.auth.Login)
			r.Post("/refresh", rt.auth.Refresh)
			r.Post("/logout", rt.auth.Logout)
		})

		r.Route("/products", func(r chi.Router) {
			r.Get("/", rt.product.List)
			r.Get("/{id}", rt.product.Get)
		})

		// Authenticated
		r.Group(func(r chi.Router) {
			r.Use(rt.authMW.Authenticate)

			r.Post("/orders", rt.order.Create)

			// Seller + Admin
			r.Group(func(r chi.Router) {
				r.Use(rt.authMW.RequireRole("seller", "admin"))
				r.Post("/seller/products", rt.product.Create)
			})
		})

		// Admin-only
		r.Group(func(r chi.Router) {
			r.Use(rt.authMW.Authenticate)
			r.Use(rt.authMW.RequireRole("admin"))

			r.Route("/admin", func(r chi.Router) {
				r.Patch("/orders/{id}/status", rt.order.UpdateStatus)
				r.Post("/products/{id}/approve", rt.product.Approve)
				r.Post("/products/{id}/reject", rt.product.Reject)
			})
		})
	})

	return r
}
