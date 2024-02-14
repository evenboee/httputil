package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/evenboee/httputil/auth"
	"github.com/evenboee/httputil/handler"
)

func main() {
	// Simple wrapper with basic auth
	withAuth := handler.Wrapper(
		auth.Basic(auth.MapChecker{"user": "pass"}),
	)

	// Create handler and add wrapper
	mux := handler.New().WithWrapperF(handler.PrintWrapper("[mux]"))

	// recovery wrapper at work
	mux.Register("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic(fmt.Errorf("%s %w", "panic", fmt.Errorf("wrapped")))
	})

	// Simple handler with auth
	mux.Register("/", withAuth(func(w http.ResponseWriter, r *http.Request) {
		u := auth.GetUsername(r)
		handler.String(w, http.StatusOK, "Hello, "+u)
	}))

	// Add handler.Func with wrapper option
	// prints [mux]
	handler.Handle(mux, "GET /people/{id}", getPerson, handler.WithWrapper(withAuth))

	// Add handler.Func with func config
	// prints [mux] and then [getPeople]
	fc := handler.NewFuncConfig().WithWrapper(handler.PrintWrapper("[getPeople]"))
	handler.HandleWith(mux, "GET /people", getPeople, fc)

	mux.MustRun("8081")
}

type getPersonParams struct {
	Path struct {
		ID string `path:"id,required"`
	}
}

// e.g. /people/123
func getPerson(w http.ResponseWriter, r *http.Request, params getPersonParams) {
	handler.String(w, http.StatusOK, "Get person: "+params.Path.ID)
}

// empty slice input gives nil slice
// slice pointer does not work
type getPeopleParams struct {
	Query struct {
		Limit      *int     `query:"limit,default=10"`
		Offset     int      `query:"offset"`
		Categories []string `query:"categories,comma"`
		IDs        []int    `query:"ids,comma"`
		Name       []string `query:"name"`
	}
	Header struct {
		Lang []string  `header:"Accept-Language,default=no en,comma"`
		Time time.Time `header:"X-Time" time_format:"2006-01-02"`
	}
}

// e.g. /people?categories=a,b,c&ids=1,2,3&name=x&name=y
func getPeople(w http.ResponseWriter, r *http.Request, params getPeopleParams) {
	handler.JSON(w, http.StatusOK, params)
}
