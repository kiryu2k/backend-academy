package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type config struct {
	Addr    string `yaml:"address"`
	Dirpath string `yaml:"files_dir_path"`
}

func LoadConfig(cfgPath string) (*config, error) {
	cfg := new(config)
	if err := cleanenv.ReadConfig(cfgPath, cfg); err != nil {
		return nil, fmt.Errorf("cannot load config: %w", err)
	}
	if cfg.Dirpath[len(cfg.Dirpath)-1] != '/' { // путь должен всегда оканчиваться на `/`
		cfg.Dirpath = cfg.Dirpath + "/"
	}
	return cfg, nil
}
