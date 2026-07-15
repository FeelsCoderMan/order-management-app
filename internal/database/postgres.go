package database

import (
	"fmt"

	"github.com/FeelsCoderMan/order-management-app/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgres() (*gorm.DB, error) {
	cfg, err := config.LoadPostgresConfig()

	if err != nil {
		return &gorm.DB{}, err
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Db)

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		return &gorm.DB{}, err
	}

	// TODO: Add auto-migration

	return db, nil
}
