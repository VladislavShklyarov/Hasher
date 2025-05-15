package server

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	lm "log-service/internal/logger/CRUD"
	"os"
)

func NewLogManager() *lm.LogManager {
	return &lm.LogManager{
		//logs: make(map[string]*gen.LogEntry),
		Loggers: map[string]*zap.Logger{
			"HTTP-service":      createLogger("http"),
			"business-service":  createLogger("business"),
			"undefined-service": createLogger("undefined"),
		},
	}
}

func createLogger(serviceName string) *zap.Logger {
	cfg := zap.NewProductionEncoderConfig()
	cfg.TimeKey = ""
	encoder := zapcore.NewJSONEncoder(cfg)

	logDir := "../log_files/"

	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	logFile, _ := os.OpenFile(logDir+"/"+serviceName+"_logs.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	core := zapcore.NewCore(encoder, zapcore.AddSync(logFile), zapcore.DebugLevel)
	return zap.New(core)
}
