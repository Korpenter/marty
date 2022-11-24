package config

import (
	"flag"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServiceAddress string `envconfig:"RUN_ADDRESS" default:"localhost:8081"`
	PostgresURI    string `envconfig:"DATABASE_URI" default:""`
	AccrualAddress string `envconfig:"ACCRUAL_SYSTEM_ADDRESS" default:"localhost:8080"`
	SecretKey      string `envconfig:"SECRET_KEY" default:"defaultKeyMARt"`
}

func NewConfig() *Config {
	var c Config
	envconfig.MustProcess("", &c)
	flag.StringVar(&c.ServiceAddress, "a", c.ServiceAddress, "адрес и порт запуска сервиса")
	flag.StringVar(&c.PostgresURI, "d", c.PostgresURI, "адрес подключения к базе данных")
	flag.StringVar(&c.AccrualAddress, "r", c.AccrualAddress, "адрес системы расчёта начислений")
	flag.Parse()
	return &c
}
