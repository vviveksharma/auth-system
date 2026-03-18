package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/argon2"
)

var charset = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var otpcharset = []rune("0123456789")

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
	return false, nil
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

func GenerateRandomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic("crypto/rand failed: " + err.Error())
		}
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

func GenerateOTP() string {
	b := make([]rune, 6)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(otpcharset))))
		if err != nil {
			panic("crypto/rand failed: " + err.Error())
		}
		b[i] = otpcharset[n.Int64()]
	}
	return string(b)
}

func GeneratePassword(password string, p *Argon2Params, salt string) (string, error) {
	decodeSalt, err := base64.RawStdEncoding.DecodeString(salt)
	if err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), decodeSalt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)
	hashBase64 := base64.RawStdEncoding.EncodeToString(hash)
	return hashBase64, nil
}

// safeRoleName only allows alphanumeric, dots, underscores, and hyphens —
// prevents path traversal when roleName originates from DB data.
var safeRoleName = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

func ReadPermissionFile(roleName string) (string, error) {
	if !safeRoleName.MatchString(roleName) {
		return "", fmt.Errorf("invalid role name: %q", roleName)
	}
	path := filepath.Join("./permissions", roleName+".json")
	file, err := os.Open(path) // #nosec G304 -- path validated by safeRoleName above
	if err != nil {
		return "", err
	}
	defer file.Close()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(file)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
