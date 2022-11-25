package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Mldlr/marty/internal/app"
	"github.com/Mldlr/marty/internal/app/constant"
	"github.com/Mldlr/marty/internal/app/container"
	"github.com/Mldlr/marty/internal/app/logging"
	"github.com/Mldlr/marty/internal/app/models"
	"github.com/Mldlr/marty/internal/app/service"
	"github.com/Mldlr/marty/internal/util/validators"
	"net/http"
)

func UserWithdraw(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			logging.Logger.Error("constant withdrawing :" + err.Error())
		}
	}()
	var withdrawal *models.Withdrawal
	if err = json.NewDecoder(r.Body).Decode(&withdrawal); err != nil {
		http.Error(w, "constant reading request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	withdrawal.Login = r.Context().Value(constant.LoginKey).(string)
	if !validators.Luhn(withdrawal.OrderID) {
		http.Error(w, "invalid order number", http.StatusUnprocessableEntity)
		return
	}
	userService := container.Container.Get("userService").(service.UserService)
	err = userService.Withdraw(r.Context(), withdrawal)
	switch {
	case errors.Is(app.ErrInsufficientBalance, err):
		http.Error(w, fmt.Sprintf("cant withdraw: %s", err), http.StatusBadRequest)
		return
	case err != nil:
		http.Error(w, fmt.Sprintf("cant withdraw: %s", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
