package config

import (
	"fmt"
	"os"
)

type Config struct {
	Env        string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	Port       string
}

const (
	EnvDevelopment = "development"
)

type ConfigLoader struct{}

func NewConfigLoader() *ConfigLoader {
	return &ConfigLoader{}
}

func (cfgLoader *ConfigLoader) Load() (Config, error) {
	env := os.Getenv("ENV")
	if env == "" {
		return Config{}, fmt.Errorf("ENV is empty")
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		return Config{}, fmt.Errorf("DB_HOST is empty")
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		return Config{}, fmt.Errorf("DB_PORT is empty")
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		return Config{}, fmt.Errorf("DB_USER is empty")
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		return Config{}, fmt.Errorf("DB_PASSWORD is empty")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		return Config{}, fmt.Errorf("DB_NAME is empty")
	}

	dbSSLMode := os.Getenv("DB_SSLMODE")
	if dbSSLMode == "" {
		return Config{}, fmt.Errorf("DB_SSLMODE is empty")
	}

	port := os.Getenv("PORT")
	if port == "" {
		return Config{}, fmt.Errorf("PORT is empty")
	}

	return Config{
		Env:        env,
		DBHost:     dbHost,
		DBPort:     dbPort,
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBName:     dbName,
		DBSSLMode:  dbSSLMode,
		Port:       port,
	}, nil
}
