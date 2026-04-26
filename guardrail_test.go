package guardrail_test

import (
	"testing"

	guardrail "github.com/vviveksharma/auth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNewGuardRail(t *testing.T) {
	// Setup in-memory database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	// Test valid configuration
	t.Run("ValidConfig", func(t *testing.T) {
		config := guardrail.Config{
			DB:        db,
			JWTSecret: "test-secret-key",
		}

		gr, err := guardrail.New(config)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if gr == nil {
			t.Error("Expected GuardRail instance, got nil")
		}
	})

	// Test missing database
	t.Run("MissingDB", func(t *testing.T) {
		config := guardrail.Config{
			JWTSecret: "test-secret-key",
		}

		_, err := guardrail.New(config)
		if err == nil {
			t.Error("Expected error for missing database, got nil")
		}
	})

	// Test missing JWT secret
	t.Run("MissingJWTSecret", func(t *testing.T) {
		config := guardrail.Config{
			DB: db,
		}

		_, err := guardrail.New(config)
		if err == nil {
			t.Error("Expected error for missing JWT secret, got nil")
		}
	})
}

func TestConfigDefaults(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	config := guardrail.Config{
		DB:        db,
		JWTSecret: "test-secret",
	}

	gr, err := guardrail.New(config)
	if err != nil {
		t.Fatalf("Failed to create GuardRail: %v", err)
	}

	if gr == nil {
		t.Fatal("GuardRail instance is nil")
	}

	// Default values should be set
	// Note: We can't access internal config directly in this pattern,
	// but we've verified it doesn't error
}

func TestHelperFunctions(t *testing.T) {
	// Test helper functions that extract values from context
	// These are simple pass-through functions, so we just verify they exist
	_ = guardrail.GetUserID
	_ = guardrail.GetRole
	_ = guardrail.GetTenantID
	_ = guardrail.GetClaims
}
