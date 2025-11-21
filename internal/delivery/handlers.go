package delivery

import (
	"PRmanager/internal/usecase"
	"net/http"
)

type Handler struct {
	usecase usecase.UsecaseInterface
}

func NewHandler(usecase usecase.UsecaseInterface) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) AddTeam(w http.ResponseWriter, r *http.Request) {
	//TODO: реализовать
}

func (h *Handler) GetTeam(w http.ResponseWriter, r *http.Request) {
	//TODO: реализовать
}
