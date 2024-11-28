package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"train_http_storage/internals/models"
)

type Files interface {
	GetAll(ctx context.Context, id int64) ([]models.File, error)
	//Add(ctx context.Context, file *models.File) error
	//GetByAlias(alias string) (models.File, error)
	//Delete(alias string) error
	//DeleteAll() error
}

type User interface {
	SingUp(ctx context.Context, input *models.SingUpInput) error
	SingIn(ctx context.Context, input *models.SingInInput) (string, error)
	ParseToken(ctx context.Context, token string) (int64, error)
}

type Handler struct {
	userService User
	fileService Files
}

func NewHandler(userService User, fileService Files) *Handler {
	return &Handler{
		userService: userService,
		fileService: fileService,
	}
}

func (h *Handler) InitRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Use(loggingMiddleware)

	router.Route("/auth", func(auth chi.Router) {
		auth.Post("/sing-up", h.singUp)
		auth.Post("/sing-in", h.singIn)
	})

	router.Route("/files", func(file chi.Router) {
		file.Use(h.authMiddleware)

		file.Get("/", h.getAll)
	})

	return router
}
