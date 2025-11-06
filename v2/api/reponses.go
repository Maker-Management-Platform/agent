package api

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, r *http.Request, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		if err := json.NewEncoder(w).Encode(v); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func Error(w http.ResponseWriter, r *http.Request, status int, err error) {
	JSON(w, r, status, map[string]string{"error": err.Error()})
}

func NoContent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
