package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"

	"RESTarchive/internals/config_parser"
	log "RESTarchive/internals/logger"
	"RESTarchive/internals/storages/postgres"
)

func main() {

	cfg := config_parser.MustLoadConfig()

	logger := log.SetupLogger(cfg.Env)

	logger.Info("LOGGER ENABLE level = %s", cfg.Env)

	// connect to database
	storagePath := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		"localhost", 5432, "postgres", "mysecretpassword")
	storage, err := postgres.NewStorage(fmt.Sprintf(storagePath))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer func() {
		err = storage.CloseStorage()
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	logger.Info("successfully connected to storage")

	logger.Info("starting server", slog.String("address", cfg.Address))

	router := httprouter.New()
	startServer(router, cfg, logger)

	logger.Error("server stopped")

}

func startServer(router *httprouter.Router, cfg *config_parser.Config, logger *slog.Logger) {
	srv := http.Server{
		Addr:         cfg.HttpConfig.Address,
		Handler:      router,
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Error("failed to start server")
	}

	logger.Error("server stopped")

}
