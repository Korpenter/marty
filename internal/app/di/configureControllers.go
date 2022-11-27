package di

import (
	"github.com/Mldlr/marty/internal/app/controllers/order"
	"github.com/Mldlr/marty/internal/app/controllers/user"
	"github.com/samber/do"
)

func configureControllers(i *do.Injector) {
	do.Provide(
		i,
		func(i *do.Injector) (*order.OrderController, error) {
			return order.NewOrderController(i), nil
		},
	)

	do.Provide(
		i,
		func(i *do.Injector) (*user.UserController, error) {
			return user.NewUserController(i), nil
		},
	)
}
