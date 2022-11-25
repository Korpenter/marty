package container

import (
	"github.com/Mldlr/marty/internal/app/config"
	"github.com/Mldlr/marty/internal/app/logging"
	"github.com/Mldlr/marty/internal/app/service"
	"github.com/Mldlr/marty/internal/app/storage"
	"github.com/sarulabs/di/v2"
)

var Container di.Container

func BuildContainer(cfg *config.Config, repo storage.Repository, userService service.UserService, orderService service.OrderService) {
	builder, _ := di.NewBuilder()
	err := builder.Add([]di.Def{
		{
			Name: "cfg",
			Build: func(ctn di.Container) (interface{}, error) {
				return cfg, nil
			},
		},
		{
			Name: "repo",
			Build: func(ctn di.Container) (interface{}, error) {
				return repo, nil
			},
		},
		{
			Name: "userService",
			Build: func(ctn di.Container) (interface{}, error) {
				return userService, nil
			},
		},
		{
			Name: "orderService",
			Build: func(ctn di.Container) (interface{}, error) {
				return orderService, nil
			},
		},
	}...)
	if err != nil {
		logging.Logger.Fatal(err.Error())
	}
	Container = builder.Build()
}
