// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package config

import (
	"os"
	"time"
)

type Config struct {
	Address           string
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	Neo4jURI          string
	Neo4jUser         string
	Neo4jPassword     string
	Neo4jDatabase     string
	Neo4jTimeout      time.Duration
}

func Load() Config {
	return Config{
		Address:           stringFromEnv("SSTPA_API_ADDR", ":8080"),
		ReadHeaderTimeout: durationFromEnv("SSTPA_API_READ_HEADER_TIMEOUT", 5*time.Second),
		WriteTimeout:      durationFromEnv("SSTPA_API_WRITE_TIMEOUT", 15*time.Second),
		Neo4jURI:          stringFromEnv("SSTPA_NEO4J_URI", ""),
		Neo4jUser:         stringFromEnv("SSTPA_NEO4J_USER", "neo4j"),
		Neo4jPassword:     stringFromEnv("SSTPA_NEO4J_PASSWORD", ""),
		Neo4jDatabase:     stringFromEnv("SSTPA_NEO4J_DATABASE", "neo4j"),
		Neo4jTimeout:      durationFromEnv("SSTPA_NEO4J_TIMEOUT", 10*time.Second),
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
