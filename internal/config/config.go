package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"sparky-back/pkg/zaplogger"
)

type Config struct {
	Server   ServerConfig     `yaml:"server"`
	Logger   zaplogger.Config `yaml:"zaplogger"`
	Database DatabaseConfig   `yaml:"database"`
}

func Load(filename string) (*Config, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}
	cfg := Config{}
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config file: %w", err)
	}

	return &cfg, nil
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}
