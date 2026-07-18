package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/FeelsCoderMan/order-management-app/internal/config"
	"github.com/FeelsCoderMan/order-management-app/internal/repository"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey int

const (
	UserKey contextKey = iota
)

type AuthMiddleware struct {
	authConfig config.AuthConfig
	logger *log.Logger
	userRepository repository.UserRepository
}

func NewAuthMiddleware(authConfig config.AuthConfig, logger *log.Logger, userRepository repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		authConfig: authConfig,
		logger: logger,
		userRepository: userRepository,
	}
}

func (m *AuthMiddleware) RequireAuthenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method used: %v", token.Method.Alg())
			}

			return []byte(m.authConfig.AccessSecret), nil
		})

		if err != nil {
			m.logger.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if token == nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userId := claims.Subject

		if userId == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if claims.ExpiresAt == nil || time.Now().After(claims.ExpiresAt.Time) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user, err := m.userRepository.GetUserById(userId)

		if err != nil {
			m.logger.Println("Failed getting user by id: ", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserKey, user)
		h.ServeHTTP(w, r.WithContext(ctx))
	})

}
