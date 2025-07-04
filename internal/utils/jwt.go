package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func CraeteJWT(userId string, roleId string, tokenType string) (string, error) {

	var expirationTime time.Time

	if tokenType == "access" {
		expirationTime = time.Now().Add(3 * time.Minute)
	} else {
		expirationTime = time.Now().Add(30 * time.Minute)
	}

	claims := jwt.MapClaims{
		"user_id": userId,
		"role_id": roleId,
		"type":    tokenType,
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
