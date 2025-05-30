package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/protobuf/types/known/durationpb"
	gen "http-service/gen"
	"http-service/internal/app"
	"http-service/internal/utils"
	"io"
	"net/http"
	"reflect"
	"time"
)

type CompositeResponse struct {
	Success            bool                 `json:"success"`
	Status             int                  `json:"status"`
	Message            string               `json:"message"`
	LogID              string               `json:"log_id,omitempty"`
	ResultID           string               `json:"result_id,omitempty"`
	LogError           string               `json:"log_error,omitempty"`
	ProcessError       string               `json:"process_error,omitempty"`
	Items              []*gen.VariableValue `json:"items,omitempty"`
	ProcessingDuration string               `json:"processing_duration"`
}

type requestJSON struct {
	Operations []operationJSON `json:"operations"`
}

type operationJSON struct {
	Type  string           `json:"type"`
	Op    string           `json:"op"`
	Var   string           `json:"var"`
	Left  utils.FlexString `json:"left"`
	Right utils.FlexString `json:"right"`
}

func ProcessDataHandler(clients *app.Clients) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		defer r.Body.Close()

		if _, err := utils.ValidateHttpRequest(r); err != nil {
			writeJSON(w, http.StatusBadRequest, CompositeResponse{
				Success: false,
				Status:  http.StatusBadRequest,
				Message: "Invalid requesst",
			})
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, CompositeResponse{
				Success: false,
				Status:  http.StatusInternalServerError,
				Message: "Failed to read request body",
			})
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		resp := CompositeResponse{
			Success: true,
			Status:  http.StatusOK,
			Message: "Request received",
		}

		var reqLogID *gen.LogID
		var logErr error

		fmt.Println(clients.LogClient)

		if !isNil(clients.LogClient) {
			fmt.Println("Мы внутри")
			reqLogID, logErr = logRequestData(r.Context(), r, clients)
			if reqLogID != nil {
				resp.LogID = reqLogID.GetId()
			}
			if logErr != nil {
				resp.LogError = logErr.Error()
				resp.Message += ", FAILED to log"
			} else {
				resp.Message += ", SUCCESSFULLY logged"
			}
		} else {
			resp.Message += ", Log service unavailable"
		}

		var resBizID *gen.LogID
		var items []*gen.VariableValue
		var procErr error
		var processingTime string
		if !isNil(clients.BusinessClient) {
			resBizID, items, processingTime, procErr = processBusinessData(r.Context(), body, clients, reqLogID)
			if resBizID != nil {
				resp.ResultID = resBizID.GetId()
			}
			if procErr != nil {
				resp.ProcessError = procErr.Error()
				resp.Message += ", FAILED processing"
			} else {
				resp.Items = items
				resp.Message += ", SUCCESSFUL processing"
				resp.ProcessingDuration = processingTime
			}
		} else {
			resp.Message += ", Business service unavailable"
		}

		if clients.LogClient == nil && clients.BusinessClient == nil {
			resp.Success = false
			resp.Status = http.StatusServiceUnavailable
			resp.Message = "Both Log and Business services unavailable"
			writeJSON(w, http.StatusServiceUnavailable, resp)
			return
		}

		writeJSON(w, http.StatusOK, resp)

	}
}

func FormatDuration(d *durationpb.Duration) string {
	// Преобразуем protobuf Duration в time.Duration
	td := d.AsDuration()

	if td < time.Microsecond {
		return fmt.Sprintf("%d ns", td.Nanoseconds())
	} else if td < time.Millisecond {
		return fmt.Sprintf("%.2f µs", float64(td.Nanoseconds())/1000)
	} else if td < time.Second {
		return fmt.Sprintf("%.2f ms", float64(td.Microseconds())/1000)
	} else {
		return fmt.Sprintf("%.2fs", td.Seconds())
	}
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch v := reflect.ValueOf(i); v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return v.IsNil()
	}
	return false
}

func logRequestData(ctx context.Context, r *http.Request, client *app.Clients) (*gen.LogID, error) {

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var reqParsed requestJSON
	if err := json.Unmarshal(bodyBytes, &reqParsed); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request body into operations: %w", err)
	}

	structured := &gen.StructuredMessage{
		Method: r.Method,
		Path:   r.URL.Path,
		Body:   make([]*gen.Operation, 0, len(reqParsed.Operations)),
	}

	for _, op := range reqParsed.Operations {

		structured.Body = append(structured.Body, &gen.Operation{
			Type:  op.Type,
			Op:    op.Op,
			Var:   op.Var,
			Left:  string(op.Left),
			Right: string(op.Right),
		})
	}

	entry := &gen.LogEntry{
		ServiceName: "HTTP-server",
		Level:       "INFO",
		Message:     structured,
		Metadata: map[string]string{
			"method": r.Method,
			"path":   r.URL.Path,
		},
		TimestampSend: time.Now().UnixMilli(),
	}

	return client.LogClient.LogDataGRPC(ctx, entry)
}

func processBusinessData(ctx context.Context, body []byte, clients *app.Clients, logID *gen.LogID) (resultID *gen.LogID,
	results []*gen.VariableValue, processingTime string, err error) {

	var reqParsed requestJSON
	if err := json.Unmarshal(body, &reqParsed); err != nil {
		return nil, nil, "", fmt.Errorf("invalid JSON: %w", err)
	}

	converted := &gen.OperationRequest{
		LogID:      logID,
		Operations: make([]*gen.Operation, 0, len(reqParsed.Operations)),
	}

	for _, op := range reqParsed.Operations {
		converted.Operations = append(converted.Operations, &gen.Operation{
			Type:  op.Type,
			Op:    op.Op,
			Var:   op.Var,
			Left:  string(op.Left),
			Right: string(op.Right),
		})
	}

	resp, err := clients.BusinessClient.Process(ctx, converted)
	results = resp.GetItems()
	if err != nil {
		return nil, nil, "", fmt.Errorf("business logic error: %w", err)
	}
	processingTime = FormatDuration(resp.GetProcessingTime())
	return resp.LogID, results, processingTime, nil
}
