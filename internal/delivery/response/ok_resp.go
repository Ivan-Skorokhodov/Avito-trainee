package response

import (
	"PRmanager/internal/models"
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

func SendOkResonseTeamCreated(team *models.TeamDTO, w http.ResponseWriter) {
	response := TeamCreatedResponse{Team: *team}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func SendOkResonseTeam(team *models.TeamDTO, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(team)
}

func SendOkResonseUser(userDto *models.UserDTO, w http.ResponseWriter) {
	response := UserResponse{User: *userDto}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func SendOkResonseReview(review *models.ReviewDTO, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(review)
}

func SendOkResonsePullRequest(pr *models.OutputCreatePullRequestDTO, w http.ResponseWriter) {
	response := CreatedPullRequestResponse{PullRequest: *pr}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func SendOKResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}
