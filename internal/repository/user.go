package repository

import (
	"errors"

	"github.com/FeelsCoderMan/order-management-app/internal/model"
	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("User not found")

type userRepository struct {
	db *gorm.DB
}

type UserRepository interface {
	Create(model.User) error
	// TODO: Handle uniqueness of username and email for register flow
	// FindByEmail will still be used for login flow
	FindByEmail(string) (*model.User, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(user model.User) error {
	return r.db.Create(&user).Error
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User;

	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}
