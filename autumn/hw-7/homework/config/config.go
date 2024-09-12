package config

import (
	"fmt"
	"homework/client"
	"homework/server"

	"github.com/ilyakaznacheev/cleanenv"
)

type config struct {
	Server *server.ServerConfig `yaml:"http_server"`
	Client *client.ClientConfig `yaml:"http_client"`
}

func New(cfgPath string) (*config, error) {
	cfg := new(config)
	if err := cleanenv.ReadConfig(cfgPath, cfg); err != nil {
		return nil, fmt.Errorf("cannot load config: %w", err)
	}
	return cfg, nil
}
