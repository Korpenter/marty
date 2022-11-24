package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Mldlr/marty/internal/app/container"
	"github.com/Mldlr/marty/internal/app/service"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
)

func UserBalance(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "jwt error", http.StatusInternalServerError)
		return
	}
	userService := container.Container.Get("userService").(service.UserService)
	balance, err := userService.GetBalance(r.Context(), claims["login"].(string))
	if err != nil {
		http.Error(w, fmt.Sprintf("cant get balance: %s", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(balance); err != nil {
		http.Error(w, "error building the response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
