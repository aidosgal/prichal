package config

import (
  "log"
  "os"
  "github.com/ilyakaznacheev/cleanenv"
  "time"
)

type Config struct {
  Env string `yaml:"env" env:"ENV" env-default:"local" env-required:"true"`
  HTTPServer `yaml:"http_server" env-required:"true"`
}

type HTTPServer struct {
  Address string `yaml:"address" env:"HTTP_SERVER_ADDRESS" env-default:"localhost:8080" env-required:"true"`
  Timeout time.Duration `yaml:"timeout" env:"HTTP_SERVER_TIMEOUT" env-required:"true"`
  IdleTimeout time.Duration `yaml:"idle_timeout" env:"HTTP_SERVER_IDLE_TIMEOUT" env-required:"true"`
}

func MustLoad () *Config {
  configPath := os.Getenv("CONFIG_PATH")
  if configPath == "" {
    log.Fatal("CONFIG_PATH is required")
  }


  // check if file exists
  if _, err := os.Stat(configPath); os.IsNotExist(err) {
    log.Fatalf("Config file not found: %s", configPath)
  }

  var cfg Config

  if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
    log.Fatalf("Failed to load config: %v", err)
  }

  return &cfg
}
