package model

import (
	"time"
)

type User struct {
	ID 			string    `gorm:"primaryKey"`
	Username    string    `gorm:"size:30;not null;uniqueIndex"`
	Email		string    `gorm:"size:255;not null;uniqueIndex"`
	Password    string    `gorm:"size:255;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
