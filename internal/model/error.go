package model

import (
	"fmt"
	"net/http"
)

const (
	InternalError             = 1999
	ValidationError           = 1001
	TransactionInProcessError = 1015
	Conflict                  = 1016
	InvalidUserError          = 1018
	ServiceUnavailableError   = 1021
	ResourceNotFoundError     = 1030
	NotImplementedError       = 1031
)

//nolint:gofumpt
var (
	//nolint:gochecknoglobals
	causeMap = map[int]causeMapping{
		InternalError:             {status: http.StatusInternalServerError, message: "Internal server error"},
		ValidationError:           {status: http.StatusBadRequest, message: "Request validation failed"},
		TransactionInProcessError: {status: http.StatusBadRequest, message: "Transaction in process"},
		Conflict:                  {status: http.StatusConflict, message: "Resource in conflict"},
		InvalidUserError:          {status: http.StatusBadRequest, message: "Invalid user"},
		ServiceUnavailableError:   {status: http.StatusBadGateway, message: "Service Unavailable"},
		ResourceNotFoundError:     {status: http.StatusNotFound, message: "Resource Not Found"},
		NotImplementedError:       {status: http.StatusNotImplemented, message: "Not Implemented"},
	}
)

type ErrorCause struct {
	Code        int64  `json:"code" example:"1030"`
	Description string `json:"description" example:"Resource Not Found"`
}

type Error struct {
	Message    string       `json:"message" example:"no rule found with key: banks_get_55603295"`
	Error      string       `json:"error" example:"Not Found"`
	Status     int          `json:"status" example:"404"`
	ErrorCause []ErrorCause `json:"cause"`
}

func NewError(causeCode int64, message string, args ...interface{}) Error {
	c := causeMap[int(causeCode)]
	cause := []ErrorCause{{Code: causeCode, Description: c.message}}
	newError := Error{
		Message:    fmt.Sprintf(message, args...),
		Status:     c.status,
		Error:      http.StatusText(c.status),
		ErrorCause: cause,
	}

	return newError
}

type causeMapping struct {
	status  int
	message string
}
