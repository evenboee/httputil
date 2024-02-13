package auth

import (
	"context"
	"net/http"
)

type contextKey string

const contextKeyAuthUsername contextKey = "username"

func setUsername(r *http.Request, user string) *http.Request {
	return r.WithContext(setUsernameC(r.Context(), user))
}

func setUsernameC(ctx context.Context, user string) context.Context {
	return context.WithValue(ctx, contextKeyAuthUsername, user)
}

func GetUsername(r *http.Request) string {
	return GetUsernameC(r.Context())
}

func GetUsernameC(ctx context.Context) string {
	if v := ctx.Value(contextKeyAuthUsername); v != nil {
		return v.(string)
	}
	return ""
}
