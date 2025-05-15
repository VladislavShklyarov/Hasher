package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJSONSuccessResponse(t *testing.T) {
	tests := []struct {
		name       string
		success    bool
		statusCode int
		message    string
		id         string
	}{
		{
			name:       "basic success",
			success:    true,
			statusCode: http.StatusOK,
			message:    "Everything is fine",
			id:         "abc123",
		},
		{
			name:       "custom message and code",
			success:    false,
			statusCode: http.StatusBadRequest,
			message:    "Invalid input",
			id:         "err-456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			writeJSONSuccessResponce(rec, tt.success, tt.statusCode, tt.message, tt.id)

			// Проверяем статус код
			if rec.Code != tt.statusCode {
				t.Errorf("expected status code %d, got %d", tt.statusCode, rec.Code)
			}

			// Проверяем заголовки
			if contentType := rec.Header().Get("Content-Type"); contentType != "application/json" {
				t.Errorf("expected Content-Type application/json, got %s", contentType)
			}

			// Проверяем тело ответа
			var resp SuccesResponse
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if resp.Success != tt.success {
				t.Errorf("expected success %v, got %v", tt.success, resp.Success)
			}
			if resp.Status != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, resp.Status)
			}
			if resp.Message != tt.message {
				t.Errorf("expected message %q, got %q", tt.message, resp.Message)
			}
			if resp.Id != tt.id {
				t.Errorf("expected id %q, got %q", tt.id, resp.Id)
			}
		})
	}
}
