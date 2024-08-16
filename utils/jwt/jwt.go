package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ut-sama-art-studio/art-market-backend/services/users"
)

var jwtSecret = []byte(getSecretKey())

func getSecretKey() string {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		return "aSecret"
	}
	return secret
}

// Generates a JWT string the user's ID for the given user
func GenerateToken(user *users.User) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().AddDate(0, 8, 0).Unix() // Token valid for 8 months

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

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
