package main

import (
	"RESTarchive/internals/config_parser"
	mylog "RESTarchive/internals/logger"
	"RESTarchive/internals/storages"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"log/slog"
	"net/http"
)

func main() {

	cfg := config_parser.MustLoadConfig()
	logger := mylog.SetupLogger(cfg.Env)
	logger.Info("LOGGER ENABLE level = %s", cfg.Env)

	mainRouter := chi.NewRouter()

	mainRouter.Use(middleware.RequestID)
	mainRouter.Use(middleware.Logger)

	// connect to database
	storagePath := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		"localhost", 5432, "postgres", "mysecretpassword")
	storage, err := postgreSQL.NewStorage(fmt.Sprintf(storagePath))
	if err != nil {
		logger.Error(err.Error())
		log.Fatal("can`t connect to storage")
	}
	logger.Info("successfully connected to storage")

	err = storage.NewUser("1233")
	if err != nil {
		logger.Error(err.Error())
		log.Fatal("can`t create user")
	}

	err = storage.UploadFiles(storages.FileToAdd{Alias: "test123",
		PathToFile: "test_path123",
		UserId:     1})
	if err != nil {
		logger.Error(err.Error())
		log.Fatal("can`t upload file")
	}

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
