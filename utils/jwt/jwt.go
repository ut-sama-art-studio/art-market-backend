package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ut-sama-art-studio/art-market-backend/services/userservice"
)

var jwtSecret = []byte(getSecretKey())

func getSecretKey() string {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		return "aSecret"
	}
	return secret
}

func GenerateKeyValueToken(key string, value string, expTime time.Time) (string, error) {
	// Generate a JWT string for a role
	claims := jwt.MapClaims{}
	claims[key] = value
	claims["exp"] = expTime.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func VerifyKeyValueToken(tokenString string, key string, value string) (jwt.MapClaims, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// Validate token and extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check expiration
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, fmt.Errorf("token is expired")
			}
		} else {
			return nil, fmt.Errorf("expiration claim is missing or invalid")
		}

		// Check key claim
		if tokenValue, ok := claims[key].(string); ok {
			if tokenValue != value {
				return nil, fmt.Errorf("role mismatch: expected %s, got %s", tokenValue, value)
			}
		} else {
			return nil, fmt.Errorf("role claim is missing or invalid")
		}

		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// Generates a JWT string the user's ID for the given user
func GenerateToken(user *userservice.User) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().AddDate(0, 4, 0).Unix() // Token valid for 4 months

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Verifies a token JWT validate
func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}

		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// Ensure the token is valid and extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check the expiration time
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, fmt.Errorf("token has expired")
			}
		} else {
			return nil, fmt.Errorf("expiration claim missing or invalid")
		}

		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
