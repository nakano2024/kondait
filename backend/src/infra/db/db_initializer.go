package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"kondait-backend/infra/config"
)

type DbInitializer struct{}

func (dbInitializr *DbInitializer) Open(cfg config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func NewDbInitializer() *DbInitializer {
	return &DbInitializer{}
}
