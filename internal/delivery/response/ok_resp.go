package response

import (
	"PRmanager/internal/models"
	"encoding/json"
	"net/http"
)

type TeamCreatedResponse struct {
	Team models.TeamDTO `json:"error"`
}

func SendOkResonseTeamCreated(team *models.TeamDTO, w http.ResponseWriter) {
	response := TeamCreatedResponse{Team: *team}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
func SendOKResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}
