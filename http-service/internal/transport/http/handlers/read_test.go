package handlers

import (
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"http-service/gen"
	"http-service/internal/app"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func makeReadLogResp(success bool, logJSON string) *gen.LogReadingResponse {
	return &gen.LogReadingResponse{
		Success: success,
		Log:     logJSON,
	}
}

func TestParseReadResponse(t *testing.T) {
	tests := []struct {
		name            string
		readLogResponse *gen.LogReadingResponse
		expect          ReadResponse
		expectErr       bool
	}{
		{
			name: "successful log parsing",
			readLogResponse: makeReadLogResp(true, `{
				"level":"info",
				"msg":"New log entry",
				"id":"abc123",
				"message":"Hello, world!",
				"source":"test-server",
				"timestamp_send":111,
				"timestamp_received":222,
				"deliveryDelayMs":"1.000 ms"
			}`),
			expect: ReadResponse{
				Success: true,
				Log: LogEntry{
					Level:           "info",
					Msg:             "New log entry",
					ID:              "abc123",
					Message:         json.RawMessage(`"Hello, world!"`),
					Source:          "test-server",
					TimestampSend:   111,
					TimestampRecv:   222,
					DeliveryDelayMs: "1.000 ms",
				},
				Error: "",
			},
			expectErr: false,
		},
		{
			name: "failed to parse invalid JSON",
			readLogResponse: makeReadLogResp(true, `{
				"level":"info",
				"msg":"New log entry",
				"id":"abc123",
				"message":"Hello, world!",
				"source":"test-server",
				"timestamp_send":111,
				"timestamp_received":222,
				"deliveryDelayMs":
			}`),
			expect:    ReadResponse{},
			expectErr: true,
		},
		{
			name: "unsuccessful response with error message",
			readLogResponse: &gen.LogReadingResponse{
				Success: false,
				Log:     `{"id":"abc123","msg":"New log entry"}`,
				Error:   "log not found",
			},
			expect: ReadResponse{
				Success: false,
				Log: LogEntry{
					Level:           "",
					Msg:             "",
					ID:              "",
					Message:         []byte("null"),
					Source:          "",
					TimestampSend:   0,
					TimestampRecv:   0,
					DeliveryDelayMs: "",
				},
				Error: "log not found",
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBytes, err := parseReadResponse(tt.readLogResponse)
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
				return
			}

			if tt.expectErr {
				return
			}

			var gotResp ReadResponse
			if err := json.Unmarshal(gotBytes, &gotResp); err != nil {
				t.Errorf("failed to unmarshal result: %v", err)
				return
			}

			if !reflect.DeepEqual(gotResp, tt.expect) {
				t.Errorf("expected response: %+v, got: %+v", tt.expect, gotResp)
			}
		})
	}
}

func TestReadLogHandler(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		mockResponse   *gen.LogReadingResponse
		mockError      error
		expectedStatus int
		expectedBody   string // for simplicity, checking a substring
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
			mockResponse:   &gen.LogReadingResponse{Success: true},
			mockError:      errors.New("internal gRPC error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "server mailfunction",
		},
		{
			name:           "log retrieval unsuccessful",
			queryParams:    "?id=123&filename=test.log",
			mockResponse:   &gen.LogReadingResponse{Success: false, Error: "log not found"},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "retrive log",
		},
		{
			name:        "successful log retrieval",
			queryParams: "?id=123&filename=test.log",
			mockResponse: &gen.LogReadingResponse{
				Success: true,
				Log: `{
	"level": "info",
	"msg": "Log message",
	"id": "123",
	"message": "This is a test log",
	"source": "http-server",
	"timestamp_send": 1620000000,
	"timestamp_received": 1620000001,
	"deliveryDelayMs": "1"
}`,
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "This is a test log",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockLogClient{
				ReadLogFunc: func(id, filename string) (*gen.LogReadingResponse, error) {
					return tt.mockResponse, tt.mockError
				},
			}

			clients := &app.Clients{
				LogClient: mockClient, // твой мок
			}

			req := httptest.NewRequest(http.MethodGet, "/read"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			handler := ReadLogHandler(clients)
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
