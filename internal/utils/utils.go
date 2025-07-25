package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/argon2"
)

type Argon2Params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

var DefaultParams = &Argon2Params{
	Memory:      64 * 1024,
	Iterations:  3,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

func CreateJWT(userId string, roleId string, tenantId string, tokenType string) (string, error) {

	var expirationTime time.Time

	if tokenType == "access" {
		expirationTime = time.Now().Add(3 * time.Minute)
	} else {
		expirationTime = time.Now().Add(30 * time.Minute)
	}

	claims := jwt.MapClaims{
		"user_id":   userId,
		"role_id":   roleId,
		"tenant_id": tenantId,
		"type":      tokenType,
		"exp":       expirationTime.Unix(),
		"is_logged": true,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func generateSalt(length uint32) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	return salt, err
}

func GeneratePasswordHash(password string, p *Argon2Params) (encodedHash string, saltBase64 string, err error) {
	salt, err := generateSalt(p.SaltLength)
	if err != nil {
		return "", "", err
	}

	hash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	saltBase64 = base64.RawStdEncoding.EncodeToString(salt)
	hashBase64 := base64.RawStdEncoding.EncodeToString(hash)

	// You can concatenate salt and hash or store separately
	return hashBase64, saltBase64, nil
}

func ComparePassword(password, storedHashBase64, storedSaltBase64 string, p *Argon2Params) (bool, error) {
	salt, err := base64.RawStdEncoding.DecodeString(storedSaltBase64)
	if err != nil {
		return false, err
	}

	newHash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)
	newHashBase64 := base64.RawStdEncoding.EncodeToString(newHash)

	// Constant-time compare
	if newHashBase64 == storedHashBase64 {
		return true, nil
	}
	return false, errors.New("password does not match")
}

func ConvertTime(input string) time.Time {
	layout := "2006-01-02" // Layout must match the input format

	parsedTime, err := time.Parse(layout, input)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return time.Time{}
	}
	return parsedTime
}
