package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	RestServerPort int `env:"REST_SERVER_PORT" env-default:"8080"`

	PostgresUser     string `env:"POSTGRES_USER" env-default:"postgres"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresDB       string `env:"POSTGRES_DB"`
	PostgresHost     string `env:"POSTGRES_HOST" env-default:"postgres"`
	PostgresPort     string `env:"POSTGRES_PORT" env-default:"5432"`
	AdminLogin       string `env:"ADMIN_LOGIN"`
	AdminPassword    string `env:"ADMIN_PASSWORD"`
}

func New() (*Config, error) {
	cfg := Config{}
	err := cleanenv.ReadConfig("./config/local.env", &cfg)
	// err := cleanenv.ReadConfig("./config/prod.env", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
