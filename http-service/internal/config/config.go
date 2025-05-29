package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	LoggerAddr   string
	BusinessAddr string
	HttpAddr     string
}

func Load() *Config {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found — assuming prod environment")
	}

	return &Config{
		LoggerAddr:   os.Getenv("LOGGER_ADDR"),
		BusinessAddr: os.Getenv("BUSINESS_ADR"),
		HttpAddr:     os.Getenv("HTTP_ADDR"),
	}
}
