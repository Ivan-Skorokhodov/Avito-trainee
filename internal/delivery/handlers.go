package delivery

import (
	"PRmanager/internal/delivery/response"
	"PRmanager/internal/models"
	"PRmanager/internal/usecase"
	"PRmanager/pkg/logs"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	appErrors "PRmanager/pkg/app_errors"
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
	var InputData models.TeamDTO
	err := json.NewDecoder(r.Body).Decode(&InputData)
	if err != nil {
		logs.PrintLog(r.Context(), "[delivery] AddTeam", err.Error())
		response.SendErrorResponse(appErrors.HttpErrParseData, w)
		return
	}

	err = h.usecase.AddTeam(r.Context(), &InputData)
	if errors.Is(err, appErrors.ErrTeamExists) {
		logs.PrintLog(r.Context(), "[delivery] AddTeam", err.Error())
		response.SendErrorResponse(appErrors.HttpErrTeamExists, w)
		return
	}

	if err != nil {
		logs.PrintLog(r.Context(), "[delivery] AddTeam", err.Error())
		response.SendErrorResponse(appErrors.HttpServerError, w)
		return
	}

	response.SendOkResonseTeamCreated(&InputData, w)
	logs.PrintLog(r.Context(), "[delivery] AddTeam", fmt.Sprintf("Team added: %+v", InputData.TeamName))
}

func (h *Handler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		logs.PrintLog(r.Context(), "[delivery] GetTeam", appErrors.ErrResourceNotFound.Error())
		response.SendErrorResponse(appErrors.HttpErrNotFound, w)
		return
	}

	team, err := h.usecase.GetTeamByName(r.Context(), teamName)
	if errors.Is(err, appErrors.ErrResourceNotFound) {
		logs.PrintLog(r.Context(), "[delivery] GetTeam", err.Error())
		response.SendErrorResponse(appErrors.HttpErrNotFound, w)
		return
	}

	if err != nil {
		logs.PrintLog(r.Context(), "[delivery] GetTeam", err.Error())
		response.SendErrorResponse(appErrors.HttpServerError, w)
		return
	}

	response.SendOkResonseTeam(team, w)
	logs.PrintLog(r.Context(), "[delivery] GetTeam", fmt.Sprintf("Team found: %+v", team.TeamName))
}
