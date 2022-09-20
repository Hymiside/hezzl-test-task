package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/Hymiside/hezzl-test-task/pkg/models"
	"github.com/Hymiside/hezzl-test-task/pkg/service"
)

type Handler struct {
	handler *chi.Mux
}

type Handlers struct {
	service *service.Service
}

func NewHandlers(s service.Service) *Handlers {
	return &Handlers{service: &s}
}

// InitHandler инициализирует хэндлеры
func (h *Handler) InitHandler(s service.Service) *chi.Mux {
	h.handler = chi.NewRouter()
	handlers := NewHandlers(s)

	h.handler.Post("/item/create", handlers.createItem)
	h.handler.Patch("/item/update", handlers.updateItem)
	h.handler.Delete("/item/remove", handlers.deleteItem)
	h.handler.Get("/item/list", handlers.getItems)

	return h.handler
}

func (s *Handlers) createItem(w http.ResponseWriter, r *http.Request) {
	var (
		ni  models.NewItem
		i   models.Item
		err error
	)

	ctx := r.Context()

	campaignId := r.URL.Query().Get("campaignId")

	if err = json.NewDecoder(r.Body).Decode(&ni); err != nil {
		ResponseError(w, "invalid request", http.StatusBadRequest)
		return
	}

	if ni.Name == "" || campaignId == "" {
		ResponseError(w, "invalid request", http.StatusBadRequest)
		return
	}

	ni.CreatedAt = time.Now()
	ni.CampaignId, err = strconv.Atoi(campaignId)
	if err != nil {
		ResponseError(w, "invalid request", http.StatusBadRequest)
		return
	}

	i, err = s.service.CreateItem(ctx, ni)
	if err != nil {
		ResponseError(w, err.Error(), http.StatusNotFound)
		return
	}
	ResponseOk(w, i)
}

func (s *Handlers) getItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	items, err := s.service.GetAllItems(ctx)
	if err != nil {
		ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResponseOk(w, items)
}

func (s *Handlers) updateItem(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		i   models.Item
	)

	ctx := r.Context()

	itemId := r.URL.Query().Get("id")
	campaignId := r.URL.Query().Get("campaignId")

	if err = json.NewDecoder(r.Body).Decode(&i); err != nil {
		ResponseError(w, "invalid request", http.StatusBadRequest)
		return
	}

	if i.Name == "" || itemId == "" || campaignId == "" {
		ResponseError(w, "invalid request", http.StatusBadRequest)
		return
	}

	i.CampaignId, err = strconv.Atoi(campaignId)
	if err != nil {
		ResponseError(w, "invalid request", http.StatusBadRequest)
		return
	}

	i.ID, err = strconv.Atoi(itemId)
	if err != nil {
		ResponseError(w, "invalid request", http.StatusBadRequest)
		return
	}

	i, err = s.service.UpdateItem(ctx, i.CampaignId, i.ID, i.Name, i.Description)
	if err != nil {
		ResponseError(w, err.Error(), http.StatusNotFound)
		return
	}
	ResponseOk(w, i)
}

func (s *Handlers) deleteItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	itemId := r.URL.Query().Get("id")
	campaignId := r.URL.Query().Get("campaignId")

	if itemId == "" || campaignId == "" {
		ResponseError(w, "invalid request", http.StatusBadRequest)
		return
	}

	campaignIdParsed, err := strconv.Atoi(campaignId)
	if err != nil {
		ResponseError(w, "invalid request", http.StatusBadRequest)
		return
	}

	itemIdParsed, err := strconv.Atoi(itemId)
	if err != nil {
		ResponseError(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err = s.service.DeleteItem(ctx, campaignIdParsed, itemIdParsed); err != nil {
		ResponseError(w, err.Error(), http.StatusNotFound)
		return
	}

	diF := &models.DeleteItem{
		ID:         itemIdParsed,
		CampaignId: campaignIdParsed,
		Removed:    true,
	}

	ResponseOk(w, diF)
}
