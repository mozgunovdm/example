package config

import (
	"os"
	"strconv"
	"strings"
)

type DatabaseConfig struct {
	Type string
	User string
	Pass string
	Host string
	Port string
	Name string
}

type HTTPServerConfig struct {
	Host string
	Port string
}

type Config struct {
	Database   DatabaseConfig
	HTTPServer HTTPServerConfig
	DebugMode  bool
}

// New returns a new Config struct
func NewConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			Type: getEnv("DB_TYPE", ""),
			User: getEnv("DB_USER", ""),
			Pass: getEnv("DB_PASS", ""),
			Host: getEnv("DB_HOST", ""),
			Port: getEnv("DB_PORT", ""),
			Name: getEnv("DB_NAME", ""),
		},
		HTTPServer: HTTPServerConfig{
			Host: getEnv("HTTP_HOST", ""),
			Port: getEnv("HTTP_PORT", ""),
		},
		DebugMode: getEnvAsBool("DEBUG_MODE", true),
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

// Helper to read an environment variable into a bool or return default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

// Helper to read an environment variable into a string slice or return default value
func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := getEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)

	return val
}
