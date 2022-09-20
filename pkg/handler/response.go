package handler

import (
	"encoding/json"
	"net/http"
)

func ResponseOk(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	resJSON, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(resJSON); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func ResponseError(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)

	resJSON, err := json.Marshal(map[string]string{"error": message})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(resJSON)
}
