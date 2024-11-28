package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	id := getIdFromRequest(r)

	files, err := h.fileService.GetAll(r.Context(), id)
	if err != nil {
		slog.Error("error with getAll handler, storage", slog.Int64("ids", id), slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	answer, err := json.Marshal(files)
	if err != nil {
		slog.Error("error with getAll handler, marshal", slog.Int64("id", id), slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(answer)
}

func getIdFromRequest(r *http.Request) int64 {
	id := r.Context().Value("userID").(int64)
	return id
}
