package utils

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGeneralValidation(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		contentType   string
		expectedError string
	}{
		{
			name:          "empty body",
			body:          "",
			contentType:   "application/json",
			expectedError: "request body is empty",
		},
		{
			name:          "invalid content type",
			body:          `{"key":"value"}`,
			contentType:   "text/plain",
			expectedError: "Content-Type must be application/json",
		},
		{
			name:          "invalid JSON",
			body:          `{"key":}`, // плохой JSON
			contentType:   "application/json",
			expectedError: "invalid JSON",
		},
		{
			name:          "valid request",
			body:          `{"key":"value"}`,
			contentType:   "application/json",
			expectedError: "", // всё ок
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", tt.contentType)

			result, err := generalValidation(req)

			if tt.expectedError != "" {
				if err == nil || !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("expected error %q, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if !bytes.Equal(result, []byte(tt.body)) {
					t.Errorf("expected body %q, got %q", tt.body, string(result))
				}

				// Проверим, что тело восстановлено корректно
				reReadBody, err := io.ReadAll(req.Body)
				if err != nil {
					t.Fatalf("failed to re-read body: %v", err)
				}
				if !bytes.Equal(reReadBody, []byte(tt.body)) {
					t.Errorf("expected restored body %q, got %q", tt.body, string(reReadBody))
				}
			}
		})
	}
}

func TestValidateHttpRequest(t *testing.T) {
	tests := []struct {
		name          string
		method        string
		body          string
		contentType   string
		expectedError string
	}{
		{
			name:          "non-POST method",
			method:        http.MethodGet,
			body:          `{"key":"value"}`,
			contentType:   "application/json",
			expectedError: "only POST requests allowed",
		},
		{
			name:          "invalid content type",
			method:        http.MethodPost,
			body:          `{"key":"value"}`,
			contentType:   "text/plain",
			expectedError: "Content-Type must be application/json",
		},
		{
			name:          "invalid JSON",
			method:        http.MethodPost,
			body:          `{"key":}`,
			contentType:   "application/json",
			expectedError: "invalid JSON",
		},
		{
			name:          "valid request",
			method:        http.MethodPost,
			body:          `{"key":"value"}`,
			contentType:   "application/json",
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", tt.contentType)

			result, err := ValidateHttpRequest(req)

			if tt.expectedError != "" {
				if err == nil || !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("expected error %q, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if !bytes.Equal(result, []byte(tt.body)) {
					t.Errorf("expected body %q, got %q", tt.body, string(result))
				}
			}
		})
	}
}
