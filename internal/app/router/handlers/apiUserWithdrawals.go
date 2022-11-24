package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Mldlr/marty/internal/app/container"
	"github.com/Mldlr/marty/internal/app/service"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
)

func UserWithdrawals(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "jwt error", http.StatusInternalServerError)
		return
	}
	userService := container.Container.Get("userService").(service.UserService)
	withdrawalsItems, err := userService.GetWithdrawals(r.Context(), claims["login"].(string))
	if err != nil {
		http.Error(w, fmt.Sprintf("cant get orders: %s", err), http.StatusInternalServerError)
		return
	}
	if len(withdrawalsItems) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(withdrawalsItems); err != nil {
		http.Error(w, "error building the response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
