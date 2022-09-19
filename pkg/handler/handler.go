package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Hymiside/hezzl-test-task/pkg/models"
	"github.com/Hymiside/hezzl-test-task/pkg/service"
	"github.com/go-chi/chi/v5"
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
	ResponseStatusOk(w, i)
}

func (s *Handlers) getItems(w http.ResponseWriter, r *http.Request) {
	items, err := s.serv.GetItems()
	if err != nil {
		ResponseError(w, err.Error(), 404)
		return
	}
	ResponseStatusOk2(w, items)
}

func (s *Handlers) updateItem(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		i   models.Item
		ui  []models.Item
	)

	itemId := r.URL.Query().Get("id")
	campaignId := r.URL.Query().Get("campaignId")

	if err = json.NewDecoder(r.Body).Decode(&i); err != nil {
		ResponseError(w, "invalid request", 404)
		return
	}

	if i.Name == "" || itemId == "" || campaignId == "" {
		ResponseError(w, "invalid request", 404)
		return
	}

	i.CampaignId, err = strconv.Atoi(campaignId)
	if err != nil {
		ResponseError(w, "invalid request", 404)
		return
	}

	i.ID, err = strconv.Atoi(itemId)
	if err != nil {
		ResponseError(w, "invalid request", 404)
		return
	}

	ui, err = s.serv.UpdateItem(i)
	if err != nil {
		ResponseError(w, err.Error(), 404)
		return
	}
	ResponseStatusOk2(w, ui)
}

func (s *Handlers) deleteItem(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		i   models.Item
		di  []models.Item
	)
	itemId := r.URL.Query().Get("id")
	campaignId := r.URL.Query().Get("campaignId")

	if itemId == "" || campaignId == "" {
		ResponseError(w, "invalid request", 404)
		return
	}

	i.CampaignId, err = strconv.Atoi(campaignId)
	if err != nil {
		ResponseError(w, "invalid request", 404)
		return
	}
	i.ID, err = strconv.Atoi(itemId)
	if err != nil {
		ResponseError(w, "invalid request", 404)
		return
	}

	di, err = s.serv.DeleteItem(i)
	if err != nil {
		ResponseError(w, err.Error(), 404)
		return
	}

	diF := &models.DeleteItem{
		ID:         di[0].ID,
		CampaignId: di[0].CampaignId,
		Removed:    di[0].Removed,
	}

	ResponseStatusOk3(w, diF)
}
