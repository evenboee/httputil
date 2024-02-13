package auth

// var (
// 	AuthorizationHeaderKey = "Authorization"
// )

// func Bearer(v Validator) handler.WrapperFunc {
// 	return func(f http.HandlerFunc) http.HandlerFunc {
// 		return func(w http.ResponseWriter, r *http.Request) {
// 			token := r.Header.Get(AuthorizationHeaderKey)
// 			if token == "" {
// 				handler.Error(w, http.StatusUnauthorized, ErrUnauthorized)
// 				return
// 			}

// 			if !strings.HasPrefix(token, "Bearer ") {
// 				handler.Error(w, http.StatusUnauthorized, ErrUnauthorized)
// 				return
// 			}

// 			token = token[7:]

// 			user, err := v.Validate(token)
// 			if err != nil {
// 				handler.Error(w, http.StatusUnauthorized, ErrUnauthorized)
// 				return
// 			}

// 			r = setUsername(r, user)
// 			f(w, r)
// 		}
// 	}
// }
