package handlers

import (
	"github.com/go-chi/chi/v5"
)

type User struct {
	SingUp(ctx context.Context, inp models)
}

type Handler struct {
}

func InitRouter(router chi.Router) *chi.Router {
	router.Use(loggingMiddleware)

}
