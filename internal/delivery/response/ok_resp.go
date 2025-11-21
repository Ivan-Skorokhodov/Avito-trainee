package response

import "net/http"

func SendOKResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func SendOkResonseCreated(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
