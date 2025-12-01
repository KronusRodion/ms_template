package config

import (
	"fmt"
	"os"
	"time"

	"github.com/goccy/go-yaml"
)


type Config struct {
    Env            string           `yaml:"env" env-default:"local"`
	GRPC           GRPCConfig       `yaml:"grpc"`
    Prometheus     PrometheusConfig `yaml:"prometheus"`
}

type GRPCConfig struct {  
    Port    *int           `yaml:"port"`  
    Timeout *time.Duration `yaml:"timeout"`  
}

type PrometheusConfig struct {  
    Port    *int           `yaml:"port"`  
}

func (cfg Config) isValid() error {
    if cfg.Env == "" {
        return fmt.Errorf("переменная env не задана в конфигурации")
    }

    err := cfg.GRPC.isValid()
    if err != nil {
        return err
    }

    return nil
}

func (g GRPCConfig) isValid() error {
    if g.Port == nil {
        return fmt.Errorf("порт gRPC не задан")
    }
    if g.Timeout == nil {
        return fmt.Errorf("таймаут для gRPC не задан")
    }
    return nil
}


func LoadConfig(path string) (*Config, error) {
    var cfg Config
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    err = yaml.Unmarshal(data, &cfg)
    if err != nil {
        return nil, err
    }

    err = cfg.isValid()

    return &cfg, err
} 