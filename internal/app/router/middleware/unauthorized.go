package middleware

import (
	"context"
	"encoding/json"
	"github.com/Mldlr/marty/internal/app/constant"
	"github.com/Mldlr/marty/internal/app/models"
	"github.com/Mldlr/marty/internal/util/validators"
	"net/http"
)

func Unauthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var cred *models.Authorization
		var err error
		if err = json.NewDecoder(r.Body).Decode(&cred); err != nil {
			http.Error(w, "error reading request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		if err = validators.ValidateAuthorization(cred); err != nil {
			http.Error(w, "data validation error", http.StatusBadRequest)
			return
		}
		req := r.WithContext(context.WithValue(r.Context(), constant.CredKey, cred))
		*r = *req
		next.ServeHTTP(w, r)
	})
}
