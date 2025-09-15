package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Host          string `yaml:"host"`
	Port          int    `yaml:"port"`
	AdminPassword string `yaml:"admin_password"`
	AdminEmail    string `yaml:"admin_email"`
	JWTSecret     string `yaml:"jwt_secret"`
}

type ClientConfig struct {
	APIBaseURL string `yaml:"api_base_url"`
	APIKey     string `yaml:"api_key"`
	TimeoutSec int    `yaml:"timeout_sec"`
}

type Config struct {
	Env    string       `yaml:"env"`
	Server ServerConfig `yaml:"server"`
	Client ClientConfig `yaml:"client"`
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	env := os.Getenv("APP_ENV")

	configPath := filepath.Join("configs", fmt.Sprintf("config.%s.yaml", env))
	data, err := os.ReadFile(configPath)
	if err != nil {
		logrus.Fatalf("Cannot read config file: %s: %v", configPath, err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		logrus.Fatalf("cannot parse config: %v", err)
	}

	return &cfg
}
