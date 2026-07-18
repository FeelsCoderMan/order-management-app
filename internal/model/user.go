package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID 			string    `gorm:"primaryKey;type:uuid"`
	Username    string    `gorm:"size:30;not null;uniqueIndex"`
	Email		string    `gorm:"size:255;not null;uniqueIndex"`
	Password    string    `gorm:"size:255;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) (error) {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}

	return nil
}
