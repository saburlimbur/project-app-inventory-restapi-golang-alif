package router

import (
	"alfdwirhmn/inventory/handler"
	appMiddleware "alfdwirhmn/inventory/middleware"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type Response struct {
	Message string `json:"message"`
}

func NewRouter(h *handler.Container, log *zap.Logger) *chi.Mux {
	r := chi.NewRouter()

	// global middleware
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// role helper
	role := appMiddleware.NewRoleMiddleware(
		h.Repositories.SessionRepo,
		h.Repositories.UserRepo,
	)

	// health check
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := Response{Message: "hello"}
		json.NewEncoder(w).Encode(resp)
	})

	r.Route("/api/v1", func(r chi.Router) {
		// public endpoint auth
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", h.Auth.Register)
			r.Post("/login", h.Auth.Login)
		})

		// admin dan super admin
		r.Route("/users", func(r chi.Router) {
			r.With(role.AllowAdmin()).Post("/", h.User.Create)
			r.With(role.AllowAdmin()).Get("/", h.User.Lists)
		})

		r.Route("/warehouse", func(r chi.Router) {
			r.With(role.AllowAdmin()).Post("/", h.Warehouse.Create)
			r.With(role.AllowRead()).Get("/", h.Warehouse.Lists)

			r.Route("/{id}", func(r chi.Router) {
				r.With(role.AllowAdmin()).Put("/", h.Warehouse.Update)
				r.With(role.AllowAdmin()).Delete("/", h.Warehouse.Delete)

				r.Route("/{id}", func(r chi.Router) {
					r.With(role.AllowAdmin()).Put("/", h.Items.Update)
				})
			})
		})

		r.Route("/racks", func(r chi.Router) {
			r.With(role.AllowAdmin()).Post("/", h.Racks.Create)
			r.With(role.AllowRead()).Get("/", h.Racks.Lists)

			r.Route("/{id}", func(r chi.Router) {
				r.With(role.AllowAdmin()).Put("/", h.Racks.Update)
				r.With(role.AllowAdmin()).Delete("/", h.Racks.Delete)
			})
		})

		r.Route("/items", func(r chi.Router) {
			r.With(role.AllowAdmin()).Post("/", h.Items.Create)
			r.With(role.AllowRead()).Get("/", h.Items.Lists)

			r.Route("/{id}", func(r chi.Router) {
				r.With(role.AllowAllRole()).Put("/", h.Items.Update)
				r.With(role.AllowAdmin()).Delete("/", h.Items.Delete)
			})
		})

		r.Route("/categories", func(r chi.Router) {
			// read all role
			r.With(role.AllowRead()).Get("/", h.Category.Lists)
			// post admin dan super admin
			r.With(role.AllowAdmin()).Post("/", h.Category.Create)

			r.Route("/{id}", func(r chi.Router) {
				r.With(role.AllowAdmin()).Put("/", h.Category.Update)
				r.With(role.AllowAdmin()).Delete("/", h.Category.Delete)
			})
		})
	})

	return r
}
