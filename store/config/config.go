package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server Server `mapstructure:"server"`

	Logging Logging `mapstructure:"logging"`

	Database Database `mapstructure:"database"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	cfg.Server.setDefaults()
	cfg.Database.setDefaults()
	cfg.Logging.setDefaults()

	if err := cfg.Server.validate(); err != nil {
		return nil, err
	}
	if err := cfg.Database.validate(); err != nil {
		return nil, err
	}
	if err := cfg.Logging.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

type Server struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func (s *Server) setDefaults() {
	if s.Host == "" {
		s.Host = "0.0.0.0"
	}
	if s.Port == 0 {
		s.Port = 8080
	}
}

func (s *Server) validate() error {
	if s.Port < 1 || s.Port > 65535 {
		return fmt.Errorf("無効なポート: %d", s.Port)
	}
	return nil
}

type Logging struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

func (l *Logging) setDefaults() {
	if l.Level == "" {
		l.Level = "info"
	}
	if l.Format == "" {
		l.Format = "json"
	}
	l.Format = strings.ToLower(l.Format)
}

func (l *Logging) validate() error {
	if l.Format != "json" && l.Format != "text" {
		return fmt.Errorf("無効なログ形式: %s (json|text)", l.Format)
	}
	return nil
}

type Database struct {
	Type     string `mapstructure:"type"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

func (d *Database) setDefaults() {
	if d.Type == "" {
		d.Type = "memory"
	}
	if d.Port == 0 {
		d.Port = 5432
	}
	if d.Type != "postgres" {
		return
	}
	if d.Host == "" {
		d.Host = "localhost"
	}
	if d.User == "" {
		d.User = "postgres"
	}
	if d.Password == "" {
		d.Password = "postgres"
	}
	if d.Name == "" {
		d.Name = "mydatabase"
	}
}

func (d *Database) validate() error {
	if d.Port < 1 || d.Port > 65535 {
		return fmt.Errorf("無効なDBポート: %d", d.Port)
	}
	return nil
}
