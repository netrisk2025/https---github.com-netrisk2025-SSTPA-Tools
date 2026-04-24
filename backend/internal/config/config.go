package config

import (
	"os"
	"time"
)

type Config struct {
	Address           string
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
}

func Load() Config {
	return Config{
		Address:           stringFromEnv("SSTPA_API_ADDR", ":8080"),
		ReadHeaderTimeout: durationFromEnv("SSTPA_API_READ_HEADER_TIMEOUT", 5*time.Second),
		WriteTimeout:      durationFromEnv("SSTPA_API_WRITE_TIMEOUT", 15*time.Second),
	}
}

func stringFromEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func durationFromEnv(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return parsed
}
