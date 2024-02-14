package bind

import (
	"net/http"
	"strings"
)

var (
	TagKeyQuery = "query"
)

func Query(obj any, r *http.Request) error {
	q := r.URL.Query()
	return BindWith(obj, func(s string) ([]string, bool) {
		v, ok := q[s]
		if !ok || v[0] == "" {
			return nil, false
		}

		return v, true
	}, TagKeyQuery)
}

var CQuerySep = ","

// CQuery is a comma separated query value
func CQuery(obj any, r *http.Request) error {
	q := r.URL.Query()
	return BindWith(obj, func(s string) ([]string, bool) {
		v := q.Get(s)
		if v == "" {
			return nil, false
		}
		return strings.Split(v, CQuerySep), true
	}, TagKeyQuery)
}
