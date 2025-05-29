package handlers

import (
	"bytes"
	"context"
	"errors"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/protobuf/types/known/durationpb"
	"http-service/gen"
	"http-service/internal/app"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestProcessDataHandler(t *testing.T) {
	tests := []struct {
		name              string
		requestBody       string
		mockLogResponse   *gen.LogID
		mockLogError      error
		mockBizResponse   *gen.OperationResponse
		mockBizError      error
		expectedStatus    int
		expectedBodyMatch []string
	}{
		{
			name:              "invalid request (malformed JSON)",
			requestBody:       `{"operations": [}`,
			expectedStatus:    http.StatusBadRequest,
			expectedBodyMatch: []string{`"success":false`, `"Invalid requesst"`},
		},
		{
			name:            "log client logs successfully, business client fails",
			requestBody:     `{"operations":[{"type":"calc","op":"add","var":"x","left":"1","right":"2"}]}`,
			mockLogResponse: &gen.LogID{Id: "log123"},
			mockBizError:    errors.New("processing error"),
			mockBizResponse: &gen.OperationResponse{},
			expectedStatus:  http.StatusOK,
			expectedBodyMatch: []string{
				`"log_id":"log123"`,
				`"process_error":"business logic error: processing error"`,
				`"message":"Request received, SUCCESSFULLY logged, FAILED processing"`,
			},
		},
		{
			name:            "both log and business succeed",
			requestBody:     `{"operations":[{"type":"calc","op":"add","var":"x","left":"1","right":"2"}]}`,
			mockLogResponse: &gen.LogID{Id: "log456"},
			mockBizResponse: &gen.OperationResponse{
				LogID:          &gen.LogID{Id: "biz789"},
				Items:          []*gen.VariableValue{{Var: "x", Value: 3}},
				ProcessingTime: durationpb.New(150 * time.Millisecond),
			},
			expectedStatus: http.StatusOK,
			expectedBodyMatch: []string{
				`"log_id":"log456"`,
				`"result_id":"biz789"`,
				`"value":3`,
				`"processing_duration":"150.00 ms"`,
				`"message":"Request received, SUCCESSFULLY logged, SUCCESSFUL processing"`,
			},
		},
		{
			name:              "both services unavailable",
			requestBody:       `{"operations":[{"type":"calc","op":"add","var":"x","left":"1","right":"2"}]}`,
			expectedStatus:    http.StatusServiceUnavailable,
			expectedBodyMatch: []string{`"success":false`, `"Both Log and Business services unavailable"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogClient := &mockLogClient{
				LogDataGRPCFunc: func(ctx context.Context, entry *gen.LogEntry) (*gen.LogID, error) {
					return tt.mockLogResponse, tt.mockLogError
				},
			}

			mockBizClient := &mockBizClient{
				ProcessFunc: func(ctx context.Context, req *gen.OperationRequest) (*gen.OperationResponse, error) {
					return tt.mockBizResponse, tt.mockBizError
				},
			}

			clients := &app.Clients{
				LogClient:      nil,
				BusinessClient: nil,
			}

			if tt.mockLogResponse != nil || tt.mockLogError != nil {
				clients.LogClient = mockLogClient
			}
			if tt.mockBizResponse != nil || tt.mockBizError != nil {
				clients.BusinessClient = mockBizClient
			}

			req := httptest.NewRequest(http.MethodPost, "/process", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			handler := ProcessDataHandler(clients)
			handler(w, req, httprouter.Params{})

			res := w.Result()
			defer res.Body.Close()
			body, _ := io.ReadAll(res.Body)

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}
			for _, expected := range tt.expectedBodyMatch {
				if !strings.Contains(string(body), expected) {
					t.Errorf("expected body to contain %q, got %q", expected, string(body))
				}
			}
		})
	}
}
