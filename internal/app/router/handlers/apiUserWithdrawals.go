package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Mldlr/marty/internal/app/container"
	"github.com/Mldlr/marty/internal/app/logging"
	"github.com/Mldlr/marty/internal/app/service"
	"net/http"
)

func UserWithdrawals(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			logging.Logger.Error("error getting withdrawals:" + err.Error())
		}
	}()
	userService := container.Container.Get("userService").(service.UserService)
	withdrawalsItems, err := userService.GetWithdrawals(r.Context())
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
