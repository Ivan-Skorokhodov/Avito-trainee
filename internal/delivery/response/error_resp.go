package response

import (
	err "PRmanager/pkg/app_errors"
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error err.HttpError `json:"error"`
}

func SendErrorResponse(httpError err.HttpError, w http.ResponseWriter) {
	response := ErrorResponse{Error: httpError}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(response)
}
