package handler

import (
	"net/http"
)

type WrapperFunc func(http.HandlerFunc) http.HandlerFunc
type WrapperHandler func(http.Handler) http.Handler

func WrapperH(w ...WrapperFunc) WrapperHandler {
	return func(h http.Handler) http.Handler {
		return Wrapper(w...)(h.ServeHTTP)
	}
}

func WrapH(h http.Handler, w ...WrapperFunc) http.Handler {
	return WrapperH(w...)(h)
}

func Wrapper(w ...WrapperFunc) WrapperFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		for i := len(w) - 1; i >= 0; i-- {
			next = w[i](next)
		}
		return next
	}
}

func DefaultWrapper(w ...WrapperFunc) WrapperFunc {
	wrappers := []WrapperFunc{
		LoggerHandler,
		DefaultRecoveryFunc,
	}
	w = append(w, wrappers...)
	return Wrapper(w...)
}

// Wrap is the same as Wrapper but for one ofs
func Wrap(f http.HandlerFunc, w ...WrapperFunc) http.HandlerFunc {
	return Wrapper(w...)(f)
}
