package handlers

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"train_http_storage/internals/models"
)

// 1. read request body 2. parsing to up/in struct 3. validate struct,
// 4. address to user service level 5. create http response

func (h *Handler) singUp(w http.ResponseWriter, r *http.Request) {
	const op = "auth.singUp"

	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("", slog.String("operation", op), slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var inp models.SingUpInput
	if err = json.Unmarshal(reqBytes, &inp); err != nil {
		slog.Error("", slog.String("operation", op), slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = inp.Validate(); err != nil {
		slog.Error("", slog.String("operation", op), slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.userService.SingUp(r.Context(), &inp)
	if err != nil {
		slog.Error("ошибка при SingUp", slog.String("operation", op), slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) singIn(w http.ResponseWriter, r *http.Request) {
	const op = "auth.singIn"

	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("", slog.String("operation", op), slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var inp models.SingInInput
	if err = json.Unmarshal(reqBytes, &inp); err != nil {
		slog.Error("", slog.String("operation", op), slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = inp.Validate(); err != nil {
		slog.Error("", slog.String("operation", op), slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := h.userService.SingIn(r.Context(), &inp)
	if err != nil {
		slog.Error("", slog.String("operation", op), slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	response, err := json.Marshal(map[string]string{
		"token": token,
	})
	if err != nil {
		slog.Error("", slog.String("operation", op), slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}
