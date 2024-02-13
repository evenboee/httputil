package main

import (
	"log"
	"net/http"

	"github.com/evenboee/httputil/auth"
	"github.com/evenboee/httputil/handler"
)

func main() {
	withAuth := handler.Wrapper(
		auth.Basic(auth.MapChecker{"user": "pass"}),
	)

	mux := http.NewServeMux()

	mux.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	})

	mux.HandleFunc("/", withAuth(func(w http.ResponseWriter, r *http.Request) {
		u := auth.GetUsername(r)
		handler.String(w, http.StatusOK, "Hello, "+u)
	}))

	log.Fatal(
		http.ListenAndServe(":8080", handler.WrapH(mux, handler.DefaultWrapper())),
	)
}
