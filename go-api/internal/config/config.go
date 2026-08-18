package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port            string
	NodeAPIURL      string
	NodeAPITimeout  time.Duration
}

func Load() (Config, error) {
	port := firstNonEmpty(os.Getenv("GO_API_PORT"), os.Getenv("PORT"), "8080")

	nodeAPIURL := os.Getenv("NODE_API_URL")
	if nodeAPIURL == "" {
		nodeAPIURL = "http://localhost:3000"
	}

	timeoutMilliseconds, err := parsePositiveInt(os.Getenv("NODE_API_TIMEOUT_MS"), 5000)
	if err != nil {
		return Config{}, fmt.Errorf("invalid NODE_API_TIMEOUT_MS: %w", err)
	}

	return Config{
		Port:           port,
		NodeAPIURL:     nodeAPIURL,
		NodeAPITimeout: time.Duration(timeoutMilliseconds) * time.Millisecond,
	}, nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}

	return ""
}

func parsePositiveInt(rawValue string, defaultValue int) (int, error) {
	if rawValue == "" {
		return defaultValue, nil
	}

	parsedValue, err := strconv.Atoi(rawValue)
	if err != nil {
		return 0, err
	}

	if parsedValue <= 0 {
		return 0, fmt.Errorf("value must be positive")
	}

	return parsedValue, nil
}
