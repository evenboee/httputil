package handler

import (
	"fmt"
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
			f := w[i]
			if f == nil {
				continue
			}
			next = f(next)
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

func DefaultWrapperH(w ...WrapperFunc) WrapperHandler {
	return WrapperH(DefaultWrapper(w...))
}

// Wrap is the same as Wrapper but for one ofs
func Wrap(f http.HandlerFunc, w ...WrapperFunc) http.HandlerFunc {
	return Wrapper(w...)(f)
}

// NopWrapper does nothing and returns the same function it received
func NopWrapper(h http.HandlerFunc) http.HandlerFunc {
	return h
}

// PrintWrapper can be used to debug wrap
func PrintWrapper(s string) WrapperFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(s)
			next(w, r)
		}
	}
}

// NopWrapperH is the same as NopWrapper but for http.Handler
func NopWrapperH(h http.Handler) http.Handler {
	return h
}
