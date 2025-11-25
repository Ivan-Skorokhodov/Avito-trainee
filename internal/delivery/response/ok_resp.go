package response

import (
	"PRmanager/internal/models"
	"PRmanager/pkg/logs"
	"context"
	"encoding/json"
	"net/http"
)

type TeamCreatedResponse struct {
	Team models.TeamDTO `json:"team"`
}

type UserResponse struct {
	User models.UserDTO `json:"user"`
}

type CreatedPullRequestResponse struct {
	PullRequest models.OutputCreatePullRequestDTO `json:"pr"`
}

type MergedPullRequestResponse struct {
	PullRequest models.OutputMergePullRequestDTO `json:"pr"`
}

type ReassignResponse struct {
	PullRequest models.OutputReassignDTO `json:"pr"`
}

func SendOkResonseTeamCreated(ctx context.Context, team *models.TeamDTO, w http.ResponseWriter) {
	response := TeamCreatedResponse{Team: *team}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logs.PrintLog(ctx, "[delivery] SendOkResonseTeamCreated", err.Error())
	}
}

func SendOkResonseTeam(ctx context.Context, team *models.TeamDTO, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(team); err != nil {
		logs.PrintLog(ctx, "[delivery] SendErrorResponse", err.Error())
	}
}

func SendOkResonseUser(ctx context.Context, userDto *models.UserDTO, w http.ResponseWriter) {
	response := UserResponse{User: *userDto}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logs.PrintLog(ctx, "[delivery] SendErrorResponse", err.Error())
	}
}

func SendOkResonseReview(ctx context.Context, review *models.ReviewDTO, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(review); err != nil {
		logs.PrintLog(ctx, "[delivery] SendErrorResponse", err.Error())
	}
}

func SendOkResonseCreatePullRequest(ctx context.Context, pr *models.OutputCreatePullRequestDTO, w http.ResponseWriter) {
	response := CreatedPullRequestResponse{PullRequest: *pr}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logs.PrintLog(ctx, "[delivery] SendErrorResponse", err.Error())
	}
}

func SendOkResonseMergePullRequest(ctx context.Context, pr *models.OutputMergePullRequestDTO, w http.ResponseWriter) {
	response := MergedPullRequestResponse{PullRequest: *pr}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logs.PrintLog(ctx, "[delivery] SendErrorResponse", err.Error())
	}
}

func SendOkResonseReassign(ctx context.Context, pr *models.OutputReassignDTO, w http.ResponseWriter) {
	response := ReassignResponse{PullRequest: *pr}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logs.PrintLog(ctx, "[delivery] SendErrorResponse", err.Error())
	}
}

func SendOKResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}
