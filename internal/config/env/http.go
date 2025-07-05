package env

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"micro-blog/internal/config"
)

type httpConfig struct {
	Port        string        `yaml:"port"  env-default:"8080"`
	Host        string        `yaml:"host"  env-default:"localhost"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func HTTPConfigLoad() (*httpConfig, error) {
	path, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	var httpCfg httpConfig

	if err = cleanenv.ReadConfig(path, &httpCfg); err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	return &httpCfg, nil
}

func (cfg *httpConfig) GetPort() string {
	return cfg.Port
}

func (cfg *httpConfig) GetHost() string {
	return cfg.Host
}

func (cfg *httpConfig) GetTimeout() time.Duration {
	return cfg.Timeout
}

func (cfg *httpConfig) GetIdleTimeout() time.Duration {
	return cfg.IdleTimeout
}
