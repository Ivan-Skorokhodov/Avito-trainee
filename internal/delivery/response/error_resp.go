package response

import (
	err "PRmanager/pkg/app_errors"
	"PRmanager/pkg/logs"
	"context"
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error err.HttpError `json:"error"`
}

func SendErrorResponse(ctx context.Context, httpError err.HttpError, w http.ResponseWriter) {
	response := ErrorResponse{Error: httpError}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpError.Status)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logs.PrintLog(ctx, "[delivery] SendErrorResponse", err.Error())
	}

}
