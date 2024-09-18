package configuration

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT  string
	HOSTS string
}

func NewConfig() (*Config, error) {
	// err := godotenv.Load("./.env")  //testing without fallback
	err := godotenv.Load() //with fallback
	if err != nil {
		return nil, err
	}
	return &Config{
		PORT:  getEnv("PORT", "8080"),
		HOSTS: getEnv("HOSTS", "127.0.0.1"), //EQUAL TO LOCALHOST
	}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
