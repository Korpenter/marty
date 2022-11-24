package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Mldlr/marty/internal/app/constant"
	"github.com/Mldlr/marty/internal/app/container"
	"github.com/Mldlr/marty/internal/app/models"
	"github.com/Mldlr/marty/internal/app/service"
	"github.com/Mldlr/marty/internal/util/validators"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
)

func UserWithdraw(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "jwt error", http.StatusInternalServerError)
		return
	}
	var withdrawal *models.Withdrawal
	if err = json.NewDecoder(r.Body).Decode(&withdrawal); err != nil {
		http.Error(w, "error reading request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	withdrawal.Login = claims["login"].(string)
	if !validators.Luhn(withdrawal.OrderID) {
		http.Error(w, "invalid order number", http.StatusUnprocessableEntity)
		return
	}
	userService := container.Container.Get("userService").(service.UserService)
	err = userService.Withdraw(r.Context(), withdrawal)
	switch {
	case errors.Is(constant.ErrInsufficientBalance, err):
		http.Error(w, fmt.Sprintf("cant withdraw: %s", err), http.StatusBadRequest)
		return
	case err != nil:
		http.Error(w, fmt.Sprintf("cant withdraw: %s", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
