package handlers

import (
	"errors"
	"fmt"
	"github.com/Mldlr/marty/internal/app"
	"github.com/Mldlr/marty/internal/app/container"
	"github.com/Mldlr/marty/internal/app/logging"
	"github.com/Mldlr/marty/internal/app/models"
	"github.com/Mldlr/marty/internal/app/service"
	"github.com/Mldlr/marty/internal/util/validators"
	"io"
	"net/http"
)

func AddOrder(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			logging.Logger.Error("constant adding order :" + err.Error())
		}
	}()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "constant reading request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	orderID := string(b)
	if !validators.Luhn(orderID) {
		http.Error(w, "invalid order number", http.StatusUnprocessableEntity)
		return
	}
	order := models.Order{
		Login:   r.Context().Value("login").(string),
		OrderID: orderID,
	}
	orderService := container.Container.Get("orderService").(service.OrderService)
	err = orderService.AddOrder(r.Context(), &order)
	switch {
	case errors.Is(app.ErrOrderAlreadyAdded, err):
		http.Error(w, fmt.Sprintf("constant adding order: %s", err), http.StatusConflict)
		return
	case errors.Is(app.ErrOrderAlreadyAddedByUser, err):
		http.Error(w, fmt.Sprintf("constant adding order: %s", err), http.StatusOK)
		return
	case err != nil:
		http.Error(w, fmt.Sprintf("constant adding order: %s", err), http.StatusInternalServerError)
		return
	}
	orderService.GetAccrual(&order)
	w.WriteHeader(http.StatusAccepted)
}
