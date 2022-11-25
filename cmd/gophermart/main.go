package main

import (
	"context"
	"fmt"
	"github.com/Mldlr/marty/internal/app/config"
	"github.com/Mldlr/marty/internal/app/container"
	"github.com/Mldlr/marty/internal/app/logging"
	"github.com/Mldlr/marty/internal/app/router"
	"github.com/Mldlr/marty/internal/app/service"
	"github.com/Mldlr/marty/internal/app/storage"
	"net/http"
)

func main() {
	defer logging.Logger.Sync()
	cfg := config.NewConfig()
	repo, err := storage.New(cfg)
	if err != nil {
		logging.Logger.Fatal(err.Error())
	}

	userService, err := service.NewUserService(cfg, repo)
	if err != nil {
		logging.Logger.Fatal(err.Error())
	}

	orderService := service.NewOrderService(cfg, repo)
	go orderService.PollAccrual()
	go orderService.UpdateOrders(context.Background())

	container.BuildContainer(cfg, repo, userService, orderService)

	r := router.NewRouter()

	logging.Logger.Info(fmt.Sprintf("Starting at: %s, with Accrual: %s", cfg.ServiceAddress, cfg.AccrualAddress))

	if err = http.ListenAndServe(cfg.ServiceAddress, r); err != nil && err != http.ErrServerClosed {
		logging.Logger.Error(err.Error())
	}
}
