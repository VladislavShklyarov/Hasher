package handlers

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	gen "http-service/gen/logger"
	"http-service/internal/app"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDeleteLogHandler(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		mockResponse   *gen.LogDeletionResponse
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "missing query params",
			queryParams:    "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "id and filename are required",
		},
		{
			name:           "gRPC call fails",
			queryParams:    "?id=123&filename=test.log",
			mockResponse:   &gen.LogDeletionResponse{Success: true, Message: "deleted successfully"}, // не важен при ошибке
			mockError:      errors.New("internal gRPC error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to delete log due to server mailfunction",
		},
		{
			name:           "log delete unsuccessful",
			queryParams:    "?id=123&filename=test.log",
			mockResponse:   &gen.LogDeletionResponse{Success: false, Message: "log not found"},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "log not found",
		},
		{
			name:           "log delete successful",
			queryParams:    "?id=123&filename=test.log",
			mockResponse:   &gen.LogDeletionResponse{Success: true, Message: "deleted successfully"},
			expectedStatus: http.StatusOK,
			expectedBody:   `"success":true`, // можно также проверить часть JSON
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockLogClient{
				DeleteLogFunc: func(id, filename string) (*gen.LogDeletionResponse, error) {
					return tt.mockResponse, tt.mockError
				},
			}

			clients := &app.Clients{
				LogClient: mockClient, // твой мок
			}

			req := httptest.NewRequest(http.MethodGet, "/delete"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			handler := DeleteLogHandler(clients)
			handler(w, req, httprouter.Params{})

			res := w.Result()
			defer res.Body.Close()
			body, _ := io.ReadAll(res.Body)

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}
			if !strings.Contains(string(body), tt.expectedBody) {
				t.Errorf("expected body to contain %q, got %q", tt.expectedBody, string(body))
			}
		})
	}
}
