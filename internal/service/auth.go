package service

import (
	"errors"
	"log"
	"time"

	"github.com/FeelsCoderMan/order-management-app/internal/config"
	"github.com/FeelsCoderMan/order-management-app/internal/dto"
	"github.com/FeelsCoderMan/order-management-app/internal/model"
	"github.com/FeelsCoderMan/order-management-app/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailAlreadyUsed   = errors.New("email already used")
	ErrInternal           = errors.New("internal server error")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type authService struct {
	authConfig  config.AuthConfig
	userRepository repository.UserRepository
	logger *log.Logger
}

type AuthService interface {
	Register(dto.CreateUserRequest) error
	Login(dto.LoginRequest) (string, time.Time, error)
}

func NewAuthService(authConfig config.AuthConfig, userRepository repository.UserRepository, logger *log.Logger) AuthService {
	return &authService{
		authConfig: authConfig,
		userRepository: userRepository,
		logger: logger,
	}
}

func (s *authService) Register(req dto.CreateUserRequest) error {
	existingUser, err := s.userRepository.GetUserByEmail(req.Email)

	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		s.logger.Println("Failed getting user by email: ", err)
		return ErrInternal
	}

	if existingUser != nil {
		return ErrEmailAlreadyUsed;
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		s.logger.Println("Failed hashing password: ", err)
		return ErrInternal
	}

	user := model.User{
		Username: req.Username,
		Email: req.Email,
		Password: string(hashedPassword),
	}

	if err := s.userRepository.Create(user); err != nil {
		s.logger.Println("Failed creating user: ", err)
		return ErrInternal
	}

	return nil
}

func (s *authService) Login(req dto.LoginRequest) (string, time.Time, error) {
	existingUser, err := s.userRepository.GetUserByEmail(req.Email)
	s.logger.Printf("User: %v", existingUser)

	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", time.Time{}, ErrInvalidCredentials
		}

		s.logger.Println("Failed retrieving user: ", err)
		return "", time.Time{}, ErrInternal
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(existingUser.Password),
		[]byte(req.Password),
	); err != nil {
		return "", time.Time{}, ErrInvalidCredentials
	}

	now := time.Now()
	expiresAt := now.Add(s.authConfig.AccessExpires)
	s.logger.Println("Existing User ID: ", existingUser.ID)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		IssuedAt : jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		Subject: existingUser.ID,
	})
	accessToken, err := token.SignedString([]byte(s.authConfig.AccessSecret))

	if err != nil {
		s.logger.Println("Failed signing token: ", err)
		return "", time.Time{}, ErrInternal
	}

	return accessToken, expiresAt, nil
}
