package config

import (
	"errors"
	"fmt"
	"os"
	"time"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `yaml:"env" env:"ENV" env-required:"true"`
	Storage Storage `yaml:"db" env-required:"true"`
	Server HTTPServer `yaml:"http_server" env-required:"true"`
	Auth Auth `yaml:"auth" env-required:"true"`
}

type Storage struct {
	Addres string `yaml:"addres" env-default:"localhost"`
	Port string `yaml:"port" env-default:"5432"`
	Database string `yaml:"database" env-required:"true"`
	User string `yaml:"user" env-default:"postgres"`
	Password string `yaml:"password" env-default:"postgres"`
	Schema string `yaml:"schema" env-default:"schema"`
}

type HTTPServer struct {
	Addres string `yaml:"addres" env-default:"localhost:8081"`
	TimeOut time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Auth struct {
	Secret string `yaml:"secret" env-required:"true"`
}

func Load() (*Config, error){
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		return nil, errors.New("CONFIG_PATH env variable not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file doesnt exist: %s", configPath)
	} 

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("cant read config: %s", err)
	} 

	return &cfg, nil
}