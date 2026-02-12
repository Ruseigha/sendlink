package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
  Environment    string
  ServerPort     string
  DatabaseURL    string
  RedisURL       string
  RedisPassword  string
  ShortURLLength int
  BaseURL        string
}

func Load() (*Config, error) {
  cfg := &Config{
    Environment:    getEnv("ENVIRONMENT", "development"),
    ServerPort:     getEnv("SERVER_PORT", "8080"),
    DatabaseURL:    getEnv("DATABASE_URL", ""),
    RedisURL:       getEnv("REDIS_URL", "localhost:6379"),
    RedisPassword:  getEnv("REDIS_PASSWORD", ""),
    ShortURLLength: getEnvAsInt("SHORT_URL_LENGTH", 6),
    BaseURL:        getEnv("BASE_URL", "http://localhost:8080"),
  }

  if cfg.DatabaseURL == "" {
      return nil, fmt.Errorf("DATABASE_URL is required")
  }

  return cfg, nil
}

func getEnv(key, defaultValue string) string {
  if value := os.Getenv(key); value != "" {
    return value
  }
  return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
  if value := os.Getenv(key); value != "" {
    if intVal, err := strconv.Atoi(value); err == nil {
      return intVal
    }
  }
  return defaultValue
}