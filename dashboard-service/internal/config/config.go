package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	WsAddr    string
	ImageDir  string
	KafkaAddr string
	BizTopic  string
	LogTopic  string
}

func Load() *Config {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found — assuming prod environment")
	}

	return &Config{
		WsAddr:    os.Getenv("WS_ADDR"),
		KafkaAddr: os.Getenv("KAFKA_ADDR"),
		ImageDir:  os.Getenv("STATIC_DIR"),
		BizTopic:  os.Getenv("BUSINESS_TOPIC"),
		LogTopic:  os.Getenv("LOG_TOPIC"),
	}
}
