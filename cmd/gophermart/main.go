package main

import (
	"github.com/Mldlr/mart/marty/internal/app/config"
	"github.com/Mldlr/mart/marty/internal/app/container"
	"github.com/Mldlr/mart/marty/internal/app/router"
	service2 "github.com/Mldlr/mart/marty/internal/app/service"
	"github.com/Mldlr/mart/marty/internal/app/storage"
	"log"
	"net/http"
)

func main() {
	cfg := config.NewConfig()
	repo, err := storage.New(cfg)
	if err != nil {
		log.Fatal("error creating repo")
	}
	userService, err := service2.NewUserService(cfg, repo)
	if err != nil {
		log.Fatal("error starting user service")
	}
	orderService := service2.NewOrderService(cfg, repo)
	go orderService.PollAccrual()
	container.BuildContainer(cfg, repo, userService, orderService)
	r := router.NewRouter()
	log.Printf("starting with cfg: %v", cfg)
	log.Fatal(http.ListenAndServe(cfg.ServiceAddress, r))
}
