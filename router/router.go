package router

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type Response struct {
	Message string `json:"message"`
}

func NewRouter(log *zap.Logger) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := Response{Message: "Hello, this is a test endpoint api!"}
		json.NewEncoder(w).Encode(resp)
	})

	return r
}
