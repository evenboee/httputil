package handler

import (
	"errors"
	"fmt"
	"net/http"
	"os"
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

type RecoveryErrFunc func(w http.ResponseWriter, r *http.Request, err error)

func RecoveryFuncE(f RecoveryErrFunc) RecoveryFunc {
	return func(w http.ResponseWriter, r *http.Request, err any) {
		e, ok := err.(error)
		if !ok {
			e = fmt.Errorf("%v", err)
		}
		f(w, r, e)
	}
}

var DefaultRecoveryFuncErrFuncOutput = os.Stdout.WriteString

var DefaultRecoveryFuncErrFunc = func(err error) {
	errChain := []error{}
	for unwrapped := err; unwrapped != nil; unwrapped = errors.Unwrap(unwrapped) {
		errChain = append(errChain, unwrapped)
	}

	DefaultRecoveryFuncErrFuncOutput("[ panic ] Recovered error\n")
	t := ""
	for _, e := range errChain {
		t += "\t"
		DefaultRecoveryFuncErrFuncOutput(t + fmt.Sprintf("%T: %s\n", e, e.Error()))
	}
}

var DefaultRecoveryFunc = RecoveryHandlerWith(
	RecoveryFuncE(func(w http.ResponseWriter, r *http.Request, err error) {
		DefaultRecoveryFuncErrFunc(err)
		Error(w, http.StatusInternalServerError, err)
	}),
)
