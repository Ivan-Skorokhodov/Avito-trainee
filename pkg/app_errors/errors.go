package apperrors

import "errors"

type HttpError struct {
	Code    string
	Message string
}

var (
	HttpErrTeamExists = HttpError{
		Code:    "TEAM_EXISTS",
		Message: "team_name already exists",
	}
	HttpServerError = HttpError{
		Code:    "SERVER_ERROR",
		Message: "server error",
	}
	HttpErrParseData = HttpError{
		Code:    "PARSE_DATA",
		Message: "can't parse data from json",
	}
	HttpErrNotFound = HttpError{
		Code:    "NOT_FOUND",
		Message: "resource team not found",
	}
)

var (
	ErrTeamExists       = errors.New("team_name already exists")
	ErrServerError      = errors.New("server error")
	ErrParseData        = errors.New("can't parse data from json")
	ErrResourceNotFound = errors.New("resource not found")
)
