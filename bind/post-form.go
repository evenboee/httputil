package bind

import (
	"net/http"
	"strings"
)

var (
	TagKeyPostForm = "form"
)

func PostForm(obj any, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	return BindWith(obj, func(s string) ([]string, bool) {
		v, ok := r.PostForm[s]
		if !ok || v[0] == "" {
			return nil, false
		}

		return v, true
	}, TagKeyPostForm)
}

var CPostFormSep = ","

func CPostForm(obj any, r *http.Request) error {
	return BindWith(obj, func(s string) ([]string, bool) {
		v, ok := r.PostForm[s]
		if !ok {
			return nil, false
		}

		if len(v) != 0 {
			v = strings.Split(v[0], CPostFormSep)
		}

		return v, true
	}, TagKeyPostForm)
}
