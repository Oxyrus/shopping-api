package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	"strings"
	"time"
)

// GenerateToken generates a JWT and signs it.
func GenerateToken(userEmail string) (string, error) {
	tokenLifeSpan, err := strconv.Atoi(GetEnv("TOKEN_LIFESPAN", "72"))
	if err != nil {
		return "", nil
	}

	claims := jwt.MapClaims{}
	claims["email"] = userEmail
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifeSpan)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenSecret := GetEnv("TOKEN_SECRET", "secret")
	return token.SignedString([]byte(tokenSecret))
}

// ValidateToken verifies if the token that comes in the Authorization header
// is valid and hasn't expired.
func ValidateToken(c *gin.Context) error {
	token := ExtractToken(c)
	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(GetEnv("TOKEN_SECRET", "secret")), nil
	})

	if err != nil {
		return err
	}

	return nil
}

// ExtractToken extracts the JWT from the Authorization header.
func ExtractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}

	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}

	return ""
}

// ExtractTokenEmail returns the email stored in the JWT that comes from
// the Authorization header.
func ExtractTokenEmail(c *gin.Context) (string, error) {
	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		email := fmt.Sprintf("%.0f", claims["email"])
		return email, nil
	}
	return "", nil
}
