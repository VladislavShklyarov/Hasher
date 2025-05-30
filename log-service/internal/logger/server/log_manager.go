package server

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"log-service/gen"
	"log-service/internal/config"
	lm "log-service/internal/logger/CRUD"
	"os"
)

func NewLogManager(cfg *config.Config) *lm.LogManager {
	return &lm.LogManager{
		Loggers: map[string]*zap.Logger{
			"HTTP-server":      createLogger("http", cfg),
			"business-server":  createLogger("business", cfg),
			"undefined-server": createLogger("undefined", cfg),
		},
		LogChanel: make(chan *gen.LogEntry, 500),
	}
}

func createLogger(serviceName string, cfg *config.Config) *zap.Logger {
	cfgZap := zap.NewProductionEncoderConfig()
	cfgZap.TimeKey = ""
	encoder := zapcore.NewJSONEncoder(cfgZap)

	logDir := cfg.LogsDir

	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	logFile, _ := os.OpenFile(logDir+"/"+serviceName+"_logs.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	fmt.Println(logDir + "/" + serviceName + "_logs.json")
	core := zapcore.NewCore(encoder, zapcore.AddSync(logFile), zapcore.DebugLevel)
	return zap.New(core)
}
