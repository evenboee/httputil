package handler

import (
	"fmt"
	"net/http"
	"reflect"
)

type ErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
	Type   string `json:"type"`
}

func NewErrorResponse(status int, err error) ErrorResponse {
	res := ErrorResponse{
		Status: status,
	}

	if err != nil {
		res.Error = err.Error()
		res.Type = reflect.TypeOf(err).String()
	}

	return res
}

func Error(w http.ResponseWriter, status int, err error) {
	JSON(w, status, NewErrorResponse(status, err))
}

type RecoveryFunc func(w http.ResponseWriter, r *http.Request, err any)

func RecoveryHandlerWith(f RecoveryFunc) WrapperFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					f(w, r, err)
				}
			}()
			next(w, r)
		}
	}
}

var DefaultRecoveryFunc = RecoveryHandlerWith(func(w http.ResponseWriter, r *http.Request, err any) {
	e, ok := err.(error)
	if !ok {
		e = fmt.Errorf("%v", err)
	}

	Error(w, http.StatusInternalServerError, e)
})
