package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"server"`

	Logging struct {
		Level  string `mapstructure:"level"`
		Format string `mapstructure:"format"`
	} `mapstructure:"logging"`

	Database struct {
		Type     string `mapstructure:"type"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
	} `mapstructure:"database"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if cfg.Server.Host == "" {
		cfg.Server.Host = "0.0.0.0"
	}
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.Logging.Level == "" {
		cfg.Logging.Level = "info"
	}
	if cfg.Logging.Format == "" {
		cfg.Logging.Format = "json"
	}
	if cfg.Database.Type == "" {
		cfg.Database.Type = "memory"
	}
	if cfg.Database.Port == 0 {
		cfg.Database.Port = 5432
	}

	if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
		return nil, fmt.Errorf("無効なポート: %d", cfg.Server.Port)
	}
	if cfg.Database.Port < 1 || cfg.Database.Port > 65535 {
		return nil, fmt.Errorf("無効なDBポート: %d", cfg.Database.Port)
	}
	cfg.Logging.Format = strings.ToLower(cfg.Logging.Format)
	if cfg.Logging.Format != "json" && cfg.Logging.Format != "text" {
		return nil, fmt.Errorf("無効なログ形式: %s (json|text)", cfg.Logging.Format)
	}

	return &cfg, nil
}
