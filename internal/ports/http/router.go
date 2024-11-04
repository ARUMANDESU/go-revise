package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (p *Port) setUpRouter() {
	if p.mux == nil {
		p.mux = chi.NewRouter()
	}
	r := p.mux

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Route("/api/v1", func(v1 chi.Router) {
		v1.Route("/users", func(r chi.Router) {
			r.Post("/register", p.handler.RegisterUser)

			r.Get("/", p.handler.GetUser)
		})

		v1.Route("/revise-items", func(r chi.Router) {
			r.Post("/", p.handler.NewReviseItem)

			r.Get("/", p.handler.GetReviseItem)
		})
	})
}
