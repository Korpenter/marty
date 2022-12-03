package main

import (
	"github.com/Mldlr/marty/internal/app/config"
	"github.com/Mldlr/marty/internal/app/di"
	"github.com/Mldlr/marty/internal/app/server"
	"go.uber.org/zap"
)

func main() {
	cfg := config.NewConfig()
	log, _ := zap.NewProduction()
	injector := di.ConfigureDependencies(cfg, log)
	log.Info("starting with cfg:",
		zap.String("Gopher Address:", cfg.ServiceAddress),
		zap.String("Accrual Address:", cfg.AccrualAddress),
	)
	srv := server.NewServer(injector)
	go srv.WaitForExitingSignal(15)
	srv.RunHTTP()
}
