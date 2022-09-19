package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Hymiside/hezzl-test-task/pkg/models"
)

func ResponseStatusOk(w http.ResponseWriter, field, message string) {
	res := make(map[string]string)
	res[field] = message
	resJSON, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, _ = w.Write(resJSON)
}

func ResponseStatusOk2(w http.ResponseWriter, message models.Item) {
	resJSON, _ := json.Marshal(message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, _ = w.Write(resJSON)
}

func ResponseStatusOk3(w http.ResponseWriter, message []models.Item) {
	res := make(map[string][]models.Item)
	res["data"] = message
	resJSON, _ := json.Marshal(message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, _ = w.Write(resJSON)
}

func ResponseError(w http.ResponseWriter, message string, code int) {
	res := make(map[string]string)
	res["message"] = message
	resJSON, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(resJSON)
}
