package config

import (
	"flag"
	"github.com/kelseyhightower/envconfig"
	"github.com/shopspring/decimal"
)

type Config struct {
	ServiceAddress string `envconfig:"RUN_ADDRESS" default:"localhost:8081"`
	PostgresURI    string `envconfig:"DATABASE_URI" default:""`
	AccrualAddress string `envconfig:"ACCRUAL_SYSTEM_ADDRESS" default:"localhost:8080"`
	SecretKey      string `envconfig:"SECRET_KEY" default:"defaultKeyMARt"`
}

func NewConfig() *Config {
	var cfg Config
	decimal.MarshalJSONWithoutQuotes = true
	envconfig.MustProcess("", &cfg)
	flag.StringVar(&cfg.ServiceAddress, "a", cfg.ServiceAddress, "адрес и порт запуска сервиса")
	flag.StringVar(&cfg.PostgresURI, "d", cfg.PostgresURI, "адрес подключения к базе данных")
	flag.StringVar(&cfg.AccrualAddress, "r", cfg.AccrualAddress, "адрес системы расчёта начислений")
	flag.Parse()
	return &cfg
}
