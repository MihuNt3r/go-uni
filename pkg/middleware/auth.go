package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthConfig struct {
	JWTSecret string
	APIKey    string
}

func AuthMiddleware(cfg AuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isValidJWT(cfg.JWTSecret, bearerTokenFromHeader(r.Header.Get("Authorization"))) || isValidAPIKey(cfg.APIKey, r.Header.Get("X-API-Key")) {
				next.ServeHTTP(w, r)
				return
			}

			writeJSONError(w, http.StatusUnauthorized, "authentication required")
		})
	}
}

func GenerateJWT(secret, subject string, ttl time.Duration) (string, error) {
	if secret == "" {
		return "", errors.New("jwt secret is empty")
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"sub": subject,
		"iat": now.Unix(),
		"exp": now.Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func isValidJWT(secret, tokenString string) bool {
	if secret == "" || tokenString == "" {
		return false
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return false
	}

	return token.Valid
}

func isValidAPIKey(expected, provided string) bool {
	if expected == "" {
		return false
	}

	return provided == expected
}

func bearerTokenFromHeader(headerValue string) string {
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(headerValue, bearerPrefix) {
		return ""
	}

	return strings.TrimSpace(strings.TrimPrefix(headerValue, bearerPrefix))
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
