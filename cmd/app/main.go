package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"os"
	"train_http_storage/internals/handlers/users"

	"train_http_storage/internals/config"
	log "train_http_storage/internals/logger"
	"train_http_storage/internals/storages/postgres"
)

func main() {

	cfg := config.MustLoadConfig()

	logger := log.SetupLogger(cfg.Env)
	slog.SetDefault(logger)

	slog.Info("enable logger", slog.String("logger level", cfg.Env))

	// connect to database
	storage, err := postgres.NewStorage(fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password))
	if err != nil {
		slog.Error("error with create NewStorage", slog.String("error", err.Error()))
		os.Exit(1)
	}

	defer func() {
		err = storage.CloseStorage()
		if err != nil {
			slog.Error("error with closing Storage", slog.String("error ", err.Error()))
		}
	}()

	slog.Info("successfully connected to storage")

	router := chi.NewRouter()

	handler := users.NewHandler()
	//	handler.Register(router)

	slog.Info("starting server", slog.String("address", cfg.Address))
	startServer(router, cfg)

}

func startServer(router chi.Router, cfg *config.Config) {
	srv := http.Server{
		Addr:         cfg.HttpConfig.Address,
		Handler:      router,
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		slog.Error("failed to start server")
	}

	slog.Error("server stopped")
}
