package handler

import (
	"encoding/json"
	"github.com/Hymiside/hezzl-test-task/pkg/models"
	"github.com/Hymiside/hezzl-test-task/pkg/service"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	handler *chi.Mux
}

type Handlers struct {
	serv *service.Service
}

func NewHandlers(s service.Service) *Handlers {
	return &Handlers{serv: &s}
}

// InitHandler функция инициализирует обработчики
func (h *Handler) InitHandler(s service.Service) *chi.Mux {
	h.handler = chi.NewRouter()
	handlers := NewHandlers(s)

	h.handler.Post("/item/create", handlers.createItem)
	h.handler.Patch("/item/update", nil)
	h.handler.Delete("/item/remove", nil)
	h.handler.Get("/item/list", handlers.getItems)

	return h.handler
}

func (s *Handlers) createItem(w http.ResponseWriter, r *http.Request) {
	var (
		ni  models.NewItem
		i   models.Item
		err error
	)

	campaignId := r.URL.Query().Get("campaignId")

	if err = json.NewDecoder(r.Body).Decode(&ni); err != nil {
		ResponseError(w, "invalid request", 404)
		return
	}

	if ni.Name == "" || ni.Description == "" || campaignId == "" {
		ResponseError(w, "invalid request", 404)
		return
	}

	ni.Removed = false
	ni.CreatedAt = time.Now()
	ni.CampaignId, err = strconv.Atoi(campaignId)
	if err != nil {
		ResponseError(w, "invalid request", 404)
		return
	}

	i, err = s.serv.CreateItem(ni)
	if err != nil {
		ResponseError(w, err.Error(), 404)
		return
	}
	ResponseStatusOk2(w, i)
}

func (s *Handlers) getItems(w http.ResponseWriter, r *http.Request) {
	items, err := s.serv.GetItems()
	if err != nil {
		ResponseError(w, err.Error(), 500)
		return
	}
	ResponseStatusOk3(w, items)
}
