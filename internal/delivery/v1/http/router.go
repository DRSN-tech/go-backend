package http

import (
	"github.com/DRSN-tech/go-backend/internal/usecase"
	"github.com/DRSN-tech/go-backend/pkg/logger"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	router *chi.Mux
	logger logger.Logger
}

func NewRouter(router *chi.Mux, logger logger.Logger) *Router {
	return &Router{router: router, logger: logger}
}

func (r *Router) Init(prUC usecase.ProductUC) {
	r.router.Route("/api/v1", func(v1 chi.Router) {
		prHandler := NewProductHandler(prUC, r.logger)
		registerProductRoutes(v1, prHandler)
	})
}

func registerProductRoutes(router chi.Router, prHandler *ProductHandler) {
	router.Route("/products", func(pr chi.Router) {
		pr.Post("/", prHandler.registerNewProduct)
	})
}
