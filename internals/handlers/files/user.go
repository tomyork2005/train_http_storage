package files

import (
	"RESTarchive/internals/handlers"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type handler struct {
}

func NewHandler() handlers.Handler {
	return &handler{}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET("api/users", h.getRegister)
}

func (h *handler) getRegister(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("register page"))
}
