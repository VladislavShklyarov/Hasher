package handlers

import (
	"encoding/json"
	gen "http-service/gen"
	"net/http"
)

type LogEntry struct {
	Level           string          `json:"level"`
	Msg             string          `json:"msg"`
	ID              string          `json:"id"`
	Message         json.RawMessage `json:"message"`
	Source          string          `json:"source"`
	TimestampSend   int64           `json:"timestamp_send"`
	TimestampRecv   int64           `json:"timestamp_received"`
	DeliveryDelayMs string          `json:"deliveryDelayMs"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
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
	return nil, nil
}

func (m *mockLogClient) LogDataGRPC(r *http.Request) (string, error) {
	if m.LogDataGRPCFunc != nil {
		return m.LogDataGRPCFunc(r)
	}
	return "", nil
}
