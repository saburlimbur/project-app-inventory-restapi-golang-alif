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

	// health check
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := Response{Message: "hello"}
		json.NewEncoder(w).Encode(resp)
	})

	r.Route("/v1", func(r chi.Router) {
		// public endpoint auth
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", h.Auth.Register)
			r.Post("/login", h.Auth.Login)
		})

		// endpoint admin dan super admin
		r.Route("/users", func(r chi.Router) {
			r.Use(
				appMiddleware.AuthMiddleware(
					h.Repositories.SessionRepo,
					h.Repositories.UserRepo,
					"admin",
					"super_admin",
				),
			)

			r.Post("/", h.User.Create)
			r.Get("/", h.User.Lists)
		})

		r.Route("/categories", func(r chi.Router) {
			r.Use(
				appMiddleware.AuthMiddleware(
					h.Repositories.SessionRepo,
					h.Repositories.UserRepo,
					"admin",
					"super_admin",
				),
			)

			r.Post("/", h.Category.Create)
		})
	})

	return r
}
