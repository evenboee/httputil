package auth

import (
	"net/http"

	"github.com/evenboee/httputil/handler"
)

func Basic(c Checker) handler.WrapperFunc {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()
			if !ok || !c.Check(user, pass) {
				handler.Error(w, http.StatusUnauthorized, ErrUnauthorized)
				return
			}

			r = setUsername(r, user)
			f(w, r)
		}
	}
}
