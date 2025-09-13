package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
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
		log.Fatalf("Cannot read config file: %s: %v", configPath, err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("cannot parse config: %v", err)
	}

	return &cfg
}
