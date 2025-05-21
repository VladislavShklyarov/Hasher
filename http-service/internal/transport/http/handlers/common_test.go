package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJSONError(t *testing.T) {
	tests := []struct {
		name       string
		success    bool
		statusCode int
		errorName  string
		reason     string
	}{
		{
			name:       "bad request error",
			success:    false,
			statusCode: http.StatusBadRequest,
			errorName:  "BadRequest",
			reason:     "Missing required field",
		},
		{
			name:       "internal error",
			success:    false,
			statusCode: http.StatusInternalServerError,
			errorName:  "InternalError",
			reason:     "Something went wrong",
		},
		{
			name:       "not found error",
			success:    false,
			statusCode: http.StatusNotFound,
			errorName:  "NotFound",
			reason:     "Resource not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			writeJSON(rec, tt.statusCode, tt.errorName+tt.reason)

			// Проверка статуса
			if rec.Code != tt.statusCode {
				t.Errorf("expected status code %d, got %d", tt.statusCode, rec.Code)
			}

			// Проверка заголовка Content-Type
			if contentType := rec.Header().Get("Content-Type"); contentType != "application/json" {
				t.Errorf("expected Content-Type application/json, got %s", contentType)
			}

			// Проверка тела ответа
			var resp ErrorResponse
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if resp.Success != tt.success {
				t.Errorf("expected success %v, got %v", tt.success, resp.Success)
			}
			if resp.Status != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, resp.Status)
			}
			if resp.Error != tt.errorName {
				t.Errorf("expected error name %q, got %q", tt.errorName, resp.Error)
			}
			if resp.Reason != tt.reason {
				t.Errorf("expected reason %q, got %q", tt.reason, resp.Reason)
			}
		})
	}
}
