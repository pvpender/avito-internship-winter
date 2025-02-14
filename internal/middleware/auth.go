package middleware

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

func Authenticator(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			token, _, err := jwtauth.FromContext(r.Context())

			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"errors": "unauthorized"}`))

				return
			}

			if token == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"errors": "unauthorized"}`))

				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(hfn)
	}
}
