package core

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

func GetErrorStatus(err error) int {
	if errors.Is(err, ENotFound) {
		slog.Error(err.Error())

		return int(ENotFound)
	}

	return http.StatusInternalServerError
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error(err.Error())
	}
}
