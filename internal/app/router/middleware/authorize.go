package middleware

import (
	"context"
	"github.com/Mldlr/marty/internal/app/constant"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"net/http"
)

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if token == nil || jwt.Validate(token) != nil || claims["login"] == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		req := r.WithContext(context.WithValue(r.Context(), constant.LoginKey, claims["login"].(string)))
		*r = *req

		next.ServeHTTP(w, r)
	})
}
