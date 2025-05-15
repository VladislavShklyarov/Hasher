package handlers

import (
	"encoding/json"
	gen "http-service/gen/logger"
	"net/http"
)

type LogEntry struct {
	Level           string `json:"level"`
	Msg             string `json:"msg"`
	ID              string `json:"id"`
	Message         string `json:"message"`
	Source          string `json:"source"`
	TimestampSend   int64  `json:"timestamp_send"`
	TimestampRecv   int64  `json:"timestamp_received"`
	DeliveryDelayMs string `json:"deliveryDelayMs"`
}

func writeJSONError(w http.ResponseWriter, success bool, statusCode int, errorName string, reason string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Success: success, Status: statusCode, Error: errorName, Reason: reason})
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
	Reason  string `json:"reason"`
}

type mockLogClient struct {
	ReadLogFunc     func(id, filename string) (*gen.LogReadingResponse, error)
	DeleteLogFunc   func(id, filename string) (*gen.LogDeletionResponse, error)
	LogDataGRPCFunc func(r *http.Request) (string, error)
}

func (m *mockLogClient) ReadLogGRPC(id, filename string) (*gen.LogReadingResponse, error) {
	return m.ReadLogFunc(id, filename)
}

func (m *mockLogClient) DeleteLogGRPC(id, filename string) (*gen.LogDeletionResponse, error) {
	if m.DeleteLogFunc != nil {
		return m.DeleteLogFunc(id, filename)
	}
	return nil, nil // или верни ошибку
}

func (m *mockLogClient) LogDataGRPC(r *http.Request) (string, error) {
	if m.LogDataGRPCFunc != nil {
		return m.LogDataGRPCFunc(r)
	}
	return "", nil // или верни ошибку
}
