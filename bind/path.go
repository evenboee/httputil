package bind

import (
	"net/http"
	"strings"
)

var (
	TagKeyPath = "path"
)

func Path(obj any, r *http.Request) error {
	return BindWith(obj, func(s string) ([]string, bool) {
		v := r.PathValue(s)
		if v == "" {
			return nil, false
		}

		return []string{v}, true
	}, TagKeyPath)
}

var CPathSep = ","

func CPath(obj any, r *http.Request) error {
	return BindWith(obj, func(s string) ([]string, bool) {
		v := r.PathValue(s)
		if v == "" {
			return nil, false
		}

		return strings.Split(v, CPathSep), true
	}, TagKeyPath)
}
