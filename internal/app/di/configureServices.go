package di

import (
	"context"
	"github.com/Mldlr/marty/internal/app/service/order"
	"github.com/Mldlr/marty/internal/app/service/user"
	"github.com/samber/do"
)

func configureServices(i *do.Injector) {

	orderService := order.NewOrderService(i)
	go orderService.PollAccrual()
	go orderService.UpdateOrders(context.Background())
	do.Provide(
		i,
		func(i *do.Injector) (order.OrderService, error) {
			return orderService, nil
		},
	)

	do.Provide(
		i,
		func(i *do.Injector) (user.UserService, error) {
			return user.NewUserService(i), nil
		},
	)
}
