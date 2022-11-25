package config

import (
	"flag"
	"github.com/kelseyhightower/envconfig"
	"github.com/shopspring/decimal"
)

type Config struct {
	ServiceAddress string `envconfig:"RUN_ADDRESS" default:"localhost:8081"`
	PostgresURI    string `envconfig:"DATABASE_URI" default:"postgresql://postgres:Gdovich6@localhost:5432/mart"`
	AccrualAddress string `envconfig:"ACCRUAL_SYSTEM_ADDRESS" default:"localhost:8080"`
	SecretKey      string `envconfig:"SECRET_KEY" default:"defaultKeyMARt"`
}

func NewConfig() *Config {
	var c Config
	decimal.MarshalJSONWithoutQuotes = true
	envconfig.MustProcess("", &c)
	flag.StringVar(&c.ServiceAddress, "a", c.ServiceAddress, "адрес и порт запуска сервиса")
	flag.StringVar(&c.PostgresURI, "d", c.PostgresURI, "адрес подключения к базе данных")
	flag.StringVar(&c.AccrualAddress, "r", c.AccrualAddress, "адрес системы расчёта начислений")
	flag.Parse()
	return &c
}
