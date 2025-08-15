package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type ConfigDB struct {
	Host      string
	Port      string
	User      string
	Password  string
	Name      string
	SSLMode   string
	JWTSecret string
}

func LoadConfig() (*ConfigDB, error) {

	err := godotenv.Load()

	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	return &ConfigDB{
		Host:      getEnv("DB_HOST", "localhost"),
		Port:      getEnv("DB_PORT", "5432"),
		User:      getEnv("DB_User", ""),
		Password:  getEnv("DB_Password", ""),
		Name:      getEnv("DB_Name", ""),
		SSLMode:   getEnv("DB_SSLMode", "disable"),
		JWTSecret: getEnv("JWT_SECRET", "default-secret-key"),
	}, nil
}

func (c *ConfigDB) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
