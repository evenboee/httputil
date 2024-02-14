package bind

import (
	"net/http"
	"strings"
)

var (
	TagKeyHeader = "header"
)

func Header(obj any, r *http.Request) error {
	return BindWith(obj, func(s string) ([]string, bool) {
		v, ok := r.Header[s]
		if !ok || v[0] == "" {
			return nil, false
		}

		return v, true
	}, TagKeyHeader)
}

var CHeaderSep = ","

func CHeader(obj any, r *http.Request) error {
	return BindWith(obj, func(s string) ([]string, bool) {
		v := r.Header.Get(s)
		if v == "" {
			return nil, false
		}

		return strings.Split(v, CHeaderSep), true
	}, TagKeyHeader)
}
