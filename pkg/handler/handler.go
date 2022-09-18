package handler

import (
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
	_ = NewHandlers(s)

	h.handler.Post("/item/create", nil)
	h.handler.Patch("/item/update", nil)
	h.handler.Delete("/item/remove", nil)
	h.handler.Get("/item/List", nil)

	return h.handler
}
