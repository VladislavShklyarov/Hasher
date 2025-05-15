package CRUD

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	gen "log-service/gen/logger"
	"log-service/internal/utils"
	"time"
)

type LogManager struct {
	gen.UnimplementedLoggerServer
	//mu      sync.RWMutex
	//logs    map[string]*gen.LogEntry
	Loggers map[string]*zap.Logger
}

func (lm *LogManager) HandleIncomingLog(ctx context.Context, entry *gen.LogEntry) (*gen.LogCreationResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	id := gen.LogID{Id: utils.GenerateID(10)}
	//lm.mu.Lock()
	//lm.logs[id.GetId()] = entry
	//lm.mu.Unlock()

	logger, ok := lm.Loggers[entry.ServiceName]
	if !ok {
		logger = lm.Loggers["undefined-service"]
	}

	WriteLogToFile(logger, entry.GetLevel(), id.GetId(), entry)

	_ = logger.Sync()

	return &gen.LogCreationResponse{Id: &id}, nil

}

func WriteLogToFile(logger *zap.Logger, level string, id string, entry *gen.LogEntry) {

	sendTs := entry.TimestampSend
	receiveTs := time.Now().UnixMilli()
	delay := float64(receiveTs - sendTs)

	formattedDelay := fmt.Sprintf("%.3f ms", delay)

	logFields := []zap.Field{
		zap.String("id", id),
		zap.String("message", entry.Message),
		zap.String("source", entry.ServiceName),
		zap.Int64("timestamp_send", sendTs),
		zap.Int64("timestamp_received", receiveTs),
		zap.String("deliveryDelayMs", formattedDelay),
	}

	switch level {
	case "debug":
		logger.Debug("New log entry", logFields...)
	case "info":
		logger.Info("New log entry", logFields...)
	case "warn":
		logger.Warn("New log entry", logFields...)
	case "error":
		logger.Error("New log entry", logFields...)
	default:
		logger.Info("New log entry", logFields...)
	}
}
