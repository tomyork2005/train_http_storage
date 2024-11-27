package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"log/slog"
	"train_http_storage/internals/models"
)

func (h *Handler) singUp(w http.ResponseWriter, r *http.Request) {
	const op = "auth.singUp"

	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("", slog.String("operation", op), slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
	}

	var inp models.SingUpInput
	if err = json.Unmarshal(reqBytes, &inp); err != nil {
		slog.Error("", slog.String("operation", op), slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
	}

	if err = inp.Validate(); err != nil {
		slog.Error("", slog.String("operation", op), slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
	}

	err = h.user

	w.WriteHeader(http.StatusOK)
}
