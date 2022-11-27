package order

import (
	"context"
	"github.com/Mldlr/marty/internal/app/models"
)

type OrderService interface {
	GetAccrual(order *models.Order)
	AddOrder(ctx context.Context, order *models.Order) error
	GetOrdersByUser(ctx context.Context) ([]models.OrderItem, error)
	UpdateOrders(ctx context.Context)
	PollAccrual()
}
