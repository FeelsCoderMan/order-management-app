package service

import (
	"errors"
	"log"

	"github.com/FeelsCoderMan/order-management-app/internal/dto"
	"github.com/FeelsCoderMan/order-management-app/internal/model"
	"github.com/FeelsCoderMan/order-management-app/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailAlreadyUsed   = errors.New("email already used")
	ErrInternal           = errors.New("internal server error")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type authService struct {
	userRepository repository.UserRepository
	logger *log.Logger
}

type AuthService interface {
	Register(dto.CreateUserRequest) error
	Login(dto.LoginRequest) (string, error)
}

func NewAuthService(userRepository repository.UserRepository, logger *log.Logger) AuthService {
	return &authService{
		userRepository: userRepository,
		logger: logger,
	}
}

func (s *authService) Register(req dto.CreateUserRequest) error {
	existingUser, err := s.userRepository.FindByEmail(req.Email)

	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		s.logger.Println("Failed checking existing user: ", err)
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

func (s *authService) Login(req dto.LoginRequest) (string, error) {
	existingUser, err := s.userRepository.FindByEmail(req.Email)

	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}

		s.logger.Println("Failed retrieving user: ", err)
		return "", ErrInternal
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(existingUser.Password),
		[]byte(req.Password),
	); err != nil {
		return "", ErrInvalidCredentials
	}

	// TODO: generate jwt token
	return "", nil
}
