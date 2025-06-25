package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func CraeteJWT(userId string, roleId string) (string, error) {

	expirationTime := time.Now().Add(3 * time.Minute)

	// Create claims
	claims := jwt.MapClaims{
		"user_id": userId,
		"role_id": roleId,
		"exp":     expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
