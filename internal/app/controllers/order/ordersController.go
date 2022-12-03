package order

import (
	"fmt"
	models2 "github.com/Mldlr/marty/internal/app/models"
	"github.com/Mldlr/marty/internal/app/service/order"
	"github.com/Mldlr/marty/internal/app/util/validators"
	"github.com/go-chi/render"
	"github.com/samber/do"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type OrderController struct {
	orderService order.OrderService
	logger       *zap.Logger
}

func NewOrderController(i *do.Injector) *OrderController {
	orderService := do.MustInvoke[order.OrderService](i)
	logger := do.MustInvoke[*zap.Logger](i)
	return &OrderController{
		orderService: orderService,
		logger:       logger,
	}
}

func (c *OrderController) HandleError(w http.ResponseWriter, r *http.Request, err error, code int) {
	c.logger.Error("request error",
		zap.String("controller", "order"),
		zap.String("url", r.URL.String()),
		zap.Error(err),
	)
	http.Error(w, err.Error(), code)
}

func (c *OrderController) AddOrder(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		c.HandleError(w, r, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	orderID := string(b)
	if !validators.Luhn(orderID) {
		c.HandleError(w, r, fmt.Errorf("invalid ID"), http.StatusUnprocessableEntity)
		return
	}
	o := models2.Order{
		Login:   r.Context().Value(models2.LoginKey{}).(string),
		OrderID: orderID,
	}
	err = c.orderService.AddOrder(r.Context(), &o)
	if err != nil {
		switch err {
		case models2.ErrOrderAlreadyAdded:
			c.HandleError(w, r, err, http.StatusConflict)
			return
		case models2.ErrOrderAlreadyAddedByUser:
			c.HandleError(w, r, err, http.StatusOK)
			return
		default:
			c.HandleError(w, r, err, http.StatusInternalServerError)
			return
		}
	}
	c.orderService.GetAccrual(&o)
	w.WriteHeader(http.StatusAccepted)
}

func (c *OrderController) OrdersByUser(w http.ResponseWriter, r *http.Request) {
	orderItems, err := c.orderService.GetOrdersByUser(r.Context())
	if err != nil {
		c.HandleError(w, r, err, http.StatusInternalServerError)
		return
	}
	if len(orderItems) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, orderItems)
}
