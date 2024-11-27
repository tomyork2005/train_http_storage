package users

import (
	"github.com/go-chi/chi/v5"
	"net/http"

	"train_http_storage/internals/handlers"
)

type handler struct {
}

func NewHandler() handlers.Handler {
	return &handler{}
}

func (h *handler) Register(router *chi.Router) {
	router.GET("/api/users", h.getRegister)
}

func (h *handler) getRegister(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("register page"))
}
