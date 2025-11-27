package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

// Config holds all configuration for the service
type Config struct {
	// Service settings
	ServiceName string `validate:"required"`
	Port        string `validate:"required"`
	Environment string `validate:"required,oneof=development staging production"`

	// Database settings
	DBHost     string `validate:"required"`
	DBPort     string `validate:"required"`
	DBUser     string `validate:"required"`
	DBPassword string `validate:"required"`
	DBName     string `validate:"required"`
	DBSSLMode  string `validate:"required,oneof=disable require verify-ca verify-full"`

	// Optional: JWT settings (for protected APIs)
	JWTSecret string

	// Optional: Allowed origins for CORS
	AllowedOrigins []string

	// Optional: Swagger host
	SwaggerHost string
}

// Load loads configuration from environment variables
func Load() *Config {
	cfg := &Config{
		// Service
		ServiceName: getEnv("SERVICE_NAME", "your-service"),
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),

		// Database
		DBHost:     getEnvRequired("DB_HOST"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnvRequired("DB_USER"),
		DBPassword: getEnvRequired("DB_PASSWORD"),
		DBName:     getEnvRequired("DB_NAME"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		// Optional
		JWTSecret:      os.Getenv("JWT_SECRET"),
		AllowedOrigins: getEnvSlice("ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
		SwaggerHost:    os.Getenv("SWAGGER_HOST"),
	}

	// Validate configuration
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		panic(fmt.Sprintf("Invalid configuration: %v", err))
	}

	return cfg
}

// DSN returns the database connection string
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

// HasJWT returns true if JWT authentication is configured
func (c *Config) HasJWT() bool {
	return c.JWTSecret != ""
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Required environment variable %s is not set", key))
	}
	return value
}

func getEnvSlice(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	// Simple comma-separated parsing
	var result []string
	start := 0
	for i := 0; i <= len(value); i++ {
		if i == len(value) || value[i] == ',' {
			if i > start {
				result = append(result, value[start:i])
			}
			start = i + 1
		}
	}
	return result
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return duration
}
