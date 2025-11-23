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
		logs.PrintLog(r.Context(), "[delivery] GetTeam", appErrors.ErrParseData.Error())
		response.SendErrorResponse(appErrors.HttpErrParseData, w)
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

func (h *Handler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	var InputData models.SetIsActiveDTO
	err := json.NewDecoder(r.Body).Decode(&InputData)
	if err != nil {
		logs.PrintLog(r.Context(), "[delivery] SetIsActive", err.Error())
		response.SendErrorResponse(appErrors.HttpErrParseData, w)
		return
	}

	userDto, err := h.usecase.SetIsActive(r.Context(), &InputData)
	if errors.Is(err, appErrors.ErrResourceNotFound) {
		logs.PrintLog(r.Context(), "[delivery] SetIsActive", err.Error())
		response.SendErrorResponse(appErrors.HttpErrNotFound, w)
		return
	}

	if err != nil {
		logs.PrintLog(r.Context(), "[delivery] SetIsActive", err.Error())
		response.SendErrorResponse(appErrors.HttpServerError, w)
		return
	}

	response.SendOkResonseUser(userDto, w)
	logs.PrintLog(r.Context(), "[delivery] SetIsActive", fmt.Sprintf("Member updated: %+v set isActive to: %+v", InputData.UserID, InputData.IsActive))
}

func (h *Handler) GetReview(w http.ResponseWriter, r *http.Request) {
	userSystemId := r.URL.Query().Get("user_id")
	if userSystemId == "" {
		logs.PrintLog(r.Context(), "[delivery] GetReview", appErrors.ErrParseData.Error())
		response.SendErrorResponse(appErrors.HttpErrParseData, w)
		return
	}

	review, err := h.usecase.GetReview(r.Context(), userSystemId)
	if errors.Is(err, appErrors.ErrResourceNotFound) {
		logs.PrintLog(r.Context(), "[delivery] GetReview", err.Error())
		response.SendErrorResponse(appErrors.HttpErrNotFound, w)
		return
	}

	if err != nil {
		logs.PrintLog(r.Context(), "[delivery] GetReview", err.Error())
		response.SendErrorResponse(appErrors.HttpServerError, w)
		return
	}

	response.SendOkResonseReview(review, w)
	logs.PrintLog(r.Context(), "[delivery] GetReview", fmt.Sprintf("Review found for user: %+v", userSystemId))
}

func (h *Handler) CreatePullRequest(w http.ResponseWriter, r *http.Request) {
	var InputData models.InputCreatePullRequestDTO
	err := json.NewDecoder(r.Body).Decode(&InputData)
	if err != nil {
		logs.PrintLog(r.Context(), "[delivery] CreatePullRequest", err.Error())
		response.SendErrorResponse(appErrors.HttpErrParseData, w)
		return
	}

	pr, err := h.usecase.CreatePullRequest(r.Context(), &InputData)
	if errors.Is(err, appErrors.ErrResourceNotFound) {
		logs.PrintLog(r.Context(), "[delivery] CreatePullRequest", err.Error())
		response.SendErrorResponse(appErrors.HttpErrNotFound, w)
		return
	}

	if errors.Is(err, appErrors.ErrPullRequestExists) {
		logs.PrintLog(r.Context(), "[delivery] CreatePullRequest", err.Error())
		response.SendErrorResponse(appErrors.HttpErrPullRequestExists, w)
		return
	}

	if err != nil {
		logs.PrintLog(r.Context(), "[delivery] CreatePullRequest", err.Error())
		response.SendErrorResponse(appErrors.HttpServerError, w)
		return
	}

	response.SendOkResonseCreatePullRequest(pr, w)
	logs.PrintLog(r.Context(), "[delivery] CreatePullRequest", fmt.Sprintf("PullRequest created: %+v", InputData.PullRequestName))
}

func (h *Handler) MergePullRequest(w http.ResponseWriter, r *http.Request) {
	var InputData models.InputMergePullRequestDTO
	err := json.NewDecoder(r.Body).Decode(&InputData)
	if err != nil {
		logs.PrintLog(r.Context(), "[delivery] MergePullRequest", err.Error())
		response.SendErrorResponse(appErrors.HttpErrParseData, w)
		return
	}

	pr, err := h.usecase.MergePullRequest(r.Context(), &InputData)
	if errors.Is(err, appErrors.ErrResourceNotFound) {
		logs.PrintLog(r.Context(), "[delivery] MergePullRequest", err.Error())
		response.SendErrorResponse(appErrors.HttpErrNotFound, w)
		return
	}

	if err != nil {
		logs.PrintLog(r.Context(), "[delivery] MergePullRequest", err.Error())
		response.SendErrorResponse(appErrors.HttpServerError, w)
		return
	}

	response.SendOkResonseMergePullRequest(pr, w)
	logs.PrintLog(r.Context(), "[delivery] MergePullRequest", fmt.Sprintf("PullRequest merged: %+v", InputData.PullRequestId))
}
