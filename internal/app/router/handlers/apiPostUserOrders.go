package handlers

import (
	"errors"
	"fmt"
	"github.com/Mldlr/marty/internal/app/constant"
	"github.com/Mldlr/marty/internal/app/container"
	"github.com/Mldlr/marty/internal/app/models"
	"github.com/Mldlr/marty/internal/app/service"
	"github.com/Mldlr/marty/internal/util/validators"
	"github.com/go-chi/jwtauth/v5"
	"io"
	"net/http"
)

func AddOrder(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "jwt error", http.StatusInternalServerError)
		return
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	orderID := string(b)
	if !validators.Luhn(orderID) {
		http.Error(w, "invalid order number", http.StatusUnprocessableEntity)
		return
	}
	order := models.Order{
		Login:   claims["login"].(string),
		OrderID: orderID,
	}
	orderService := container.Container.Get("orderService").(service.OrderService)
	err = orderService.AddOrder(r.Context(), &order)
	switch {
	case errors.Is(constant.ErrOrderAlreadyAdded, err):
		http.Error(w, fmt.Sprintf("error adding order: %s", err), http.StatusConflict)
		return
	case errors.Is(constant.ErrOrderAlreadyAddedByUser, err):
		http.Error(w, fmt.Sprintf("error adding order: %s", err), http.StatusOK)
		return
	case err != nil:
		http.Error(w, fmt.Sprintf("error adding order: %s", err), http.StatusInternalServerError)
		return
	}
	orderService.GetAccrual(&order)
	w.WriteHeader(http.StatusAccepted)
}
