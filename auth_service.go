package guardrail

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
	"gorm.io/gorm"
)

// AuthService provides registration, login, and user management functionality
type AuthService struct {
	gr *GuardRail
}

// NewAuthService creates a new auth service instance
func (gr *GuardRail) NewAuthService() *AuthService {
	return &AuthService{gr: gr}
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`      // Optional, defaults to "user"
	TenantID  string `json:"tenant_id"` // Required if multi-tenant is enabled
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role"`      // Optional, for RBAC systems
	TenantID string `json:"tenant_id"` // Required if multi-tenant is enabled
}

// AuthResponse represents the response after login/registration
type AuthResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	UserID       string    `json:"user_id"`
	Role         string    `json:"role,omitempty"`
	TenantID     string    `json:"tenant_id,omitempty"`
}

// User represents a user in the database
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Password  string    `gorm:"not null"` // Hashed password
	Salt      string    `gorm:"not null"`
	FirstName string
	LastName  string
	Role      string    `gorm:"default:'user'"`
	TenantID  uuid.UUID `gorm:"type:uuid;index"`
	IsActive  bool      `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

// Register creates a new user account
func (as *AuthService) Register(req RegisterRequest) (*AuthResponse, error) {
	// Validate tenant_id if multi-tenant is enabled
	if as.gr.config.EnableMultiTenant && req.TenantID == "" {
		return nil, fmt.Errorf("tenant_id is required")
	}

	// Set default role if not provided
	if req.Role == "" {
		req.Role = "user"
	}

	// Check if user already exists
	var existingUser User
	err := as.gr.db.Where("email = ?", req.Email).First(&existingUser).Error
	if err == nil {
		return nil, fmt.Errorf("user with this email already exists")
	} else if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Hash password
	salt := generateSalt()
	hashedPassword := hashPassword(req.Password, salt)

	// Create user
	user := User{
		ID:        uuid.New(),
		Email:     req.Email,
		Password:  hashedPassword,
		Salt:      salt,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
		IsActive:  true,
	}

	if req.TenantID != "" {
		tenantUUID, err := uuid.Parse(req.TenantID)
		if err != nil {
			return nil, fmt.Errorf("invalid tenant_id: %w", err)
		}
		user.TenantID = tenantUUID
	}

	// Save to database
	if err := as.gr.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens
	return as.generateAuthResponse(user)
}

// Login authenticates a user and returns tokens
func (as *AuthService) Login(req LoginRequest) (*AuthResponse, error) {
	// Validate tenant_id if multi-tenant is enabled
	if as.gr.config.EnableMultiTenant && req.TenantID == "" {
		return nil, fmt.Errorf("tenant_id is required")
	}

	// Find user by email
	var user User
	query := as.gr.db.Where("email = ? AND is_active = true", req.Email)

	if req.TenantID != "" {
		tenantUUID, err := uuid.Parse(req.TenantID)
		if err != nil {
			return nil, fmt.Errorf("invalid tenant_id: %w", err)
		}
		query = query.Where("tenant_id = ?", tenantUUID)
	}

	err := query.First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invalid email or password")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Verify password
	if !verifyPassword(req.Password, user.Password, user.Salt) {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Check role if RBAC is enabled and role is specified
	if as.gr.config.EnableRBAC && req.Role != "" && user.Role != req.Role {
		return nil, fmt.Errorf("invalid role for this user")
	}

	// Generate tokens
	return as.generateAuthResponse(user)
}

// RefreshToken generates a new access token from a refresh token
func (as *AuthService) RefreshToken(refreshToken string) (*AuthResponse, error) {
	// Verify the refresh token
	claims, err := as.gr.verifyJWT(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Extract user_id from claims
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id in token: %w", err)
	}

	// Fetch user from database
	var user User
	if err := as.gr.db.Where("id = ? AND is_active = true", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found or inactive: %w", err)
	}

	// Generate new tokens
	return as.generateAuthResponse(user)
}

// Logout invalidates a token by adding it to the blacklist
func (as *AuthService) Logout(token string) error {
	if as.gr.redis == nil {
		return fmt.Errorf("Redis is required for logout functionality")
	}

	// Add token to blacklist
	ctx := as.gr.db.Statement.Context
	return as.gr.redis.Set(ctx, "blacklist:"+token, "logged_out", 24*time.Hour).Err()
}

// generateAuthResponse creates tokens and returns auth response
func (as *AuthService) generateAuthResponse(user User) (*AuthResponse, error) {
	now := time.Now()
	accessExpiry := now.Add(as.gr.config.AccessTokenExpiry)
	refreshExpiry := now.Add(as.gr.config.RefreshTokenExpiry)

	// Create access token claims
	accessClaims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"role":    user.Role,
		"exp":     accessExpiry.Unix(),
		"iat":     now.Unix(),
		"type":    "access",
	}

	if as.gr.config.EnableMultiTenant {
		accessClaims["tenant_id"] = user.TenantID.String()
	}

	// Create refresh token claims
	refreshClaims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     refreshExpiry.Unix(),
		"iat":     now.Unix(),
		"type":    "refresh",
	}

	// Generate tokens
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessTokenString, err := accessToken.SignedString(as.gr.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	refreshTokenString, err := refreshToken.SignedString(as.gr.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	response := &AuthResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresAt:    accessExpiry,
		UserID:       user.ID.String(),
		Role:         user.Role,
	}

	if as.gr.config.EnableMultiTenant {
		response.TenantID = user.TenantID.String()
	}

	return response, nil
}

// Password hashing functions using Argon2
const (
	saltLength = 16
)

func generateSalt() string {
	salt := make([]byte, saltLength)
	rand.Read(salt)
	return base64.RawStdEncoding.EncodeToString(salt)
}

func hashPassword(password, salt string) string {
	saltBytes, _ := base64.RawStdEncoding.DecodeString(salt)
	hash := argon2.IDKey([]byte(password), saltBytes, 1, 64*1024, 4, 32)
	return base64.RawStdEncoding.EncodeToString(hash)
}

func verifyPassword(password, hashedPassword, salt string) bool {
	newHash := hashPassword(password, salt)
	return subtle.ConstantTimeCompare([]byte(newHash), []byte(hashedPassword)) == 1
}
