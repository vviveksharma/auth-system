package guardrail

import (
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Config holds all the configuration needed for GuardRail middleware
type Config struct {
	// Database connection (required)
	DB *gorm.DB

	// JWT Secret for token signing/verification (required)
	JWTSecret string

	// Token expiration durations (optional, defaults provided)
	AccessTokenExpiry  time.Duration // Default: 15 minutes
	RefreshTokenExpiry time.Duration // Default: 7 days

	// Redis client for caching (optional but recommended)
	RedisClient *redis.Client

	// Custom error messages (optional)
	ErrorMessages ErrorMessages

	// Enable/disable features
	EnableRBAC        bool // Enable Role-Based Access Control (default: true)
	EnableMultiTenant bool // Enable multi-tenant support (default: false)
}

// ErrorMessages allows customization of error responses
type ErrorMessages struct {
	Unauthorized  string
	Forbidden     string
	InvalidToken  string
	MissingToken  string
	ExpiredToken  string
	InvalidAppKey string
	MissingAppKey string
	InternalError string
}

// setDefaults sets default values for optional configuration
func (c *Config) setDefaults() {
	if c.AccessTokenExpiry == 0 {
		c.AccessTokenExpiry = 15 * time.Minute
	}
	if c.RefreshTokenExpiry == 0 {
		c.RefreshTokenExpiry = 7 * 24 * time.Hour
	}

	// Set default error messages if not provided
	if c.ErrorMessages.Unauthorized == "" {
		c.ErrorMessages.Unauthorized = "Unauthorized access"
	}
	if c.ErrorMessages.Forbidden == "" {
		c.ErrorMessages.Forbidden = "Access forbidden"
	}
	if c.ErrorMessages.InvalidToken == "" {
		c.ErrorMessages.InvalidToken = "Invalid or expired token"
	}
	if c.ErrorMessages.MissingToken == "" {
		c.ErrorMessages.MissingToken = "Authorization header is required"
	}
	if c.ErrorMessages.ExpiredToken == "" {
		c.ErrorMessages.ExpiredToken = "Token has expired. Please login again"
	}
	if c.ErrorMessages.InvalidAppKey == "" {
		c.ErrorMessages.InvalidAppKey = "Invalid application key"
	}
	if c.ErrorMessages.MissingAppKey == "" {
		c.ErrorMessages.MissingAppKey = "Application key is required"
	}
	if c.ErrorMessages.InternalError == "" {
		c.ErrorMessages.InternalError = "Internal server error"
	}
}

// Validate checks if required configuration is provided
func (c *Config) Validate() error {
	if c.DB == nil {
		return &ConfigError{Field: "DB", Message: "database connection is required"}
	}
	if c.JWTSecret == "" {
		return &ConfigError{Field: "JWTSecret", Message: "JWT secret is required"}
	}
	return nil
}

// ConfigError represents a configuration validation error
type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return "guardrail config error: " + e.Field + " - " + e.Message
}
