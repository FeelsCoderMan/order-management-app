package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/FeelsCoderMan/order-management-app/internal/repository"
	"github.com/FeelsCoderMan/order-management-app/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey int

const (
	userKey contextKey = iota
)

type AuthMiddleware struct {
	logger *log.Logger
	userRepository repository.UserRepository
}

func NewAuthMiddleware(logger *log.Logger, userRepository repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		logger: logger,
		userRepository: userRepository,
	}
}

func (m *AuthMiddleware) RequireAuthenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authCookie, err := r.Cookie("Authorization")

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(authCookie.String(), func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method used: %v", token.Method.Alg())
			}

			return []byte("test"), nil
		})

		if err != nil {
			m.logger.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if token == nil || !token.Valid {
			utils.WriteJSON(w, http.StatusUnauthorized, &UnauthorizedResponse{})
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.RegisteredClaims)

		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userId := claims.Subject

		if userId != "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if claims.ExpiresAt == nil || time.Now().After(claims.ExpiresAt.Time) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user, err := m.userRepository.GetUserById(userId)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, user)
		h.ServeHTTP(w, r.WithContext(ctx))
	})

}
