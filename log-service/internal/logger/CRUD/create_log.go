package CRUD

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"log-service/gen"
	"log-service/internal/utils"
	"time"
)

type LogManager struct {
	gen.UnimplementedLoggerServer
	Loggers   map[string]*zap.Logger
	LogChanel chan *gen.LogEntry
}

func (lm *LogManager) HandleIncomingLog(ctx context.Context, entry *gen.LogEntry) (*gen.LogCreationResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	id := gen.LogID{Id: utils.GenerateID(10)}
	logger, ok := lm.Loggers[entry.ServiceName]
	if !ok {
		logger = lm.Loggers["undefined-server"]
	}

	WriteLogToFile(logger, entry.GetLevel(), id.GetId(), entry, lm.LogChanel)

	_ = logger.Sync()

	return &gen.LogCreationResponse{Id: &id}, nil

}

func WriteLogToFile(logger *zap.Logger, level string, id string, entry *gen.LogEntry, logChan chan *gen.LogEntry) {

	sendTs := entry.TimestampSend
	receiveTs := time.Now().UnixMilli()
	delay := float64(receiveTs - sendTs)

	formattedDelay := fmt.Sprintf("%.3f ms", delay)

	msgJSON, err := protojson.Marshal(entry.Message)
	if err != nil {
		msgJSON = []byte(fmt.Sprintf(`"error serializing message: %v"`, err))
	}

	logFields := []zap.Field{
		zap.String("id", id),
		zap.Any("message", json.RawMessage(msgJSON)),
		zap.String("path", entry.Message.GetPath()),
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

	if entry.ServiceName == "business-server" {
		select {
		case logChan <- entry:
			fmt.Printf("Log %s successfully send to chanel\n", id)
		default:
			log.Println("LogChanel is full, dropping message")
		}
	}
}
