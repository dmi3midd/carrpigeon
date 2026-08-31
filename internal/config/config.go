package config

import (
	"fmt"
	"os"
	"time"

	"go.yaml.in/yaml/v3"
)

type Database struct {
	Name         string        `yaml:"name"`
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	User         string        `yaml:"user"`
	Password     string        `yaml:"password"`
	SSLMode      string        `yaml:"sslmode"`
	MaxOpenConns int           `yaml:"maxOpenConns"`
	MaxIdleConns int           `yaml:"maxIdleConns"`
	MaxIdleTime  time.Duration `yaml:"maxIdleTime"`
}

type HTTPServer struct {
	Address      string        `yaml:"address"`
	IdleTimeout  time.Duration `yaml:"idleTimeout"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
	CORS         CORS          `yaml:"cors"`
}

type CORS struct {
	AllowedOrigins []string `yaml:"allowedOrigins"`
}

type SMTP struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type Logger struct {
	Level string `yaml:"level"`
}

type Config struct {
	Database   `yaml:"database"`
	HTTPServer `yaml:"httpServer"`
	SMTP       `yaml:"smtp"`
	Logger     `yaml:"logger"`
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return cfg, nil
}
