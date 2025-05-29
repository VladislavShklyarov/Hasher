package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	LoggerAddr  string
	LogsDir     string
	KafkaBroker string
	KafkaTopic  string
}

func Load() *Config {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found — assuming prod environment")
	}

	return &Config{
		LoggerAddr:  os.Getenv("LOGGER_ADDR"),
		LogsDir:     os.Getenv("LOGS_DIR"),
		KafkaBroker: os.Getenv("KAFKA_BROKER"),
		KafkaTopic:  os.Getenv("KAFKA_TOPIC"),
	}
}
