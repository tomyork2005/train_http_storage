package main

import (
	"RESTarchive/internals/config_parser"
	mylog "RESTarchive/internals/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
)

func main() {

	cfg := config_parser.MustLoadConfig()
	logger := mylog.SetupLogger("dev")
	logger.Info("LOGGER ENABLE level = %s", cfg.Env)

	mainRouter := chi.NewRouter()

	mainRouter.Use(middleware.RequestID)
	mainRouter.Use(middleware.Logger)

	mainRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	logger.Info("starting server", slog.String("address", cfg.Address))

	srv := http.Server{
		Addr:         cfg.HttpConfig.Address,
		Handler:      mainRouter,
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Error("failed to start server")
	}

	logger.Error("server stopped")

}
