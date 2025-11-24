package apperrors

import (
	"errors"
	"net/http"
)

type HttpError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

var (
	HttpErrTeamExists = HttpError{
		Code:    "TEAM_EXISTS",
		Message: "team_name already exists",
		Status:  http.StatusBadRequest,
	}
	HttpServerError = HttpError{
		Code:    "SERVER_ERROR",
		Message: "server error",
		Status:  http.StatusInternalServerError,
	}
	HttpErrParseData = HttpError{
		Code:    "PARSE_DATA",
		Message: "can't parse data from json",
		Status:  http.StatusBadRequest,
	}
	HttpErrNotFound = HttpError{
		Code:    "NOT_FOUND",
		Message: "resource not found",
		Status:  http.StatusNotFound,
	}
	HttpErrPullRequestExists = HttpError{
		Code:    "PR_EXISTS",
		Message: "PR id already exists",
		Status:  http.StatusConflict,
	}
	HttpErrPullRequestMerged = HttpError{
		Code:    "PR_MERGED",
		Message: "cannot reassign on merged PR",
		Status:  http.StatusConflict,
	}
)

var (
	ErrTeamExists        = errors.New("team_name already exists")
	ErrServerError       = errors.New("server error")
	ErrParseData         = errors.New("can't parse data from json")
	ErrResourceNotFound  = errors.New("resource not found")
	ErrPullRequestExists = errors.New("pr id already exists")
	ErrPullRequestMerged = errors.New("cannot reassign on merged PR")
)
