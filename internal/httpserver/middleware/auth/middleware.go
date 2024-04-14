package auth

import (
	"context"
	"net/http"
)

type RequiredMiddleware struct {
	Next http.Handler
}

func (rm *RequiredMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	if token != "admin" && token != "user" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ctx := context.WithValue(r.Context(), "token", token)
	rm.Next.ServeHTTP(w, r.WithContext(ctx))
}

func NewRequiredMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return &RequiredMiddleware{Next: next}
	}
}

func RoleRequired(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("token")

			if token != role {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), "token", token)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
