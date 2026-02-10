package config

import (
	"fmt"

	sharedconfig "github.com/adrputra/face-recognition-svc/shared/config"
)

const DefaultConfigPath = "/run/secrets/presence-config"
const DefaultGRPCPort = 50051

type Config struct {
	GRPC GRPCConfig `yaml:"grpc"`
}

type GRPCConfig struct {
	Port int `yaml:"port"`
}

// Example YAML:
//
//	grpc:
//	  port: 50051
func Load(path string) (Config, error) {
	if path == "" {
		path = DefaultConfigPath
	}

	var cfg Config
	if err := sharedconfig.LoadYAML(path, &cfg); err != nil {
		return Config{}, fmt.Errorf("load config file: %w", err)
	}

	if cfg.GRPC.Port == 0 {
		cfg.GRPC.Port = DefaultGRPCPort
	}

	return cfg, nil
}
