package http

import (
	_ "github.com/DRSN-tech/go-backend/docs" // Импорт сгенерированных файлов
	"github.com/DRSN-tech/go-backend/internal/usecase"
	"github.com/DRSN-tech/go-backend/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Router struct {
	router *chi.Mux
	logger logger.Logger
}

func NewRouter(router *chi.Mux, logger logger.Logger) *Router {
	return &Router{router: router, logger: logger}
}

func (r *Router) Init(prUC usecase.ProductUC) {
	r.router.Use(middleware.Logger)    // Пишет логи запросов в консоль
	r.router.Use(middleware.Recoverer) // Не дает серверу упасть при панике
	
	r.router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // ссылка на JSON
	))
	
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
