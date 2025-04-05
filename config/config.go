package config

import (
	"os"
)

func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

var (
	Port        = GetEnv("GRPC_PORT", "9001")
	GithubToken = os.Getenv("GITHUB_TOKEN")
)
