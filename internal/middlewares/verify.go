package middlewares

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/vviveksharma/auth/cache"
)

func VerifyJWT(tokenStr string) (jwt.MapClaims, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	if len(secret) == 0 {
		return nil, fmt.Errorf("JWT_SECRET is not set")
	}

	// 1. Check if token is blacklisted (expired/invalid)
	blacklisted, _ := cache.Exists("blacklist:" + tokenStr)
	if blacklisted {
		return nil, fmt.Errorf("token is expired or invalid")
	}

	// 2. Check cache for valid tokens
	var cachedClaims jwt.MapClaims
	if err := cache.Get("token:"+tokenStr, &cachedClaims); err == nil {
		return cachedClaims, nil
	}

	// 3. Cache miss - validate token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		// Token is invalid/expired - blacklist it for 1 hour
		cache.Set("blacklist:"+tokenStr, "expired", 0)
		return nil, fmt.Errorf("error while verifying token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Cache valid token until its expiry
		if exp, ok := claims["exp"].(float64); ok {
			expiresAt := time.Unix(int64(exp), 0)
			ttl := time.Until(expiresAt)
			if ttl > 0 {
				cache.Set("token:"+tokenStr, claims, ttl)
			}
		}
		return claims, nil
	}

	// Token is invalid - blacklist it
	cache.Set("blacklist:"+tokenStr, "invalid", 1*time.Hour)
	return nil, fmt.Errorf("invalid token or claims")
}
