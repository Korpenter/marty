package container

import (
	"github.com/Mldlr/mart/marty/internal/app/config"
	service2 "github.com/Mldlr/mart/marty/internal/app/service"
	"github.com/Mldlr/mart/marty/internal/app/storage"
	"github.com/sarulabs/di/v2"
	"log"
)

var Container di.Container

func BuildContainer(cfg *config.Config, repo storage.Repository, userService service2.UserService, orderService service2.OrderService) {
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
		log.Fatalln(err)
	}
	Container = builder.Build()
}
