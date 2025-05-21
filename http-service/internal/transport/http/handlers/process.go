package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	gen "http-service/gen"
	"http-service/internal/app"
	"http-service/internal/utils"
	"io"
	"net/http"
	"time"
)

type CompositeResponse struct {
	Success      bool                 `json:"success"`
	Status       int                  `json:"status"`
	Message      string               `json:"message"`
	LogID        string               `json:"log_id,omitempty"`
	ResultID     string               `json:"result_id,omitempty"`
	LogError     string               `json:"log_error,omitempty"`
	ProcessError string               `json:"process_error,omitempty"`
	Items        []*gen.VariableValue `json:"items,omitempty"`
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

		reqLogID, logErr := logRequestData(r.Context(), r, clients)
		resID, items, procErr := processBusinessData(r.Context(), body, clients, reqLogID)

		fmt.Println("Log from logRequestData: " + reqLogID.GetId())
		fmt.Println("Log from processBusinessData: " + resID.GetId())

		// Отправляем ID логирования и бизнес серверу тоже -- чтобы затем можно было отследить операции, на которые сделан запрос
		resp := CompositeResponse{
			Success:  true,
			Status:   http.StatusOK,
			Message:  "Request processed",
			LogID:    reqLogID.GetId(),
			ResultID: resID.GetId(),
		}

		if logErr != nil {
			resp.LogError = logErr.Error()
			resp.Message += ", FAILED to log"
		} else {
			resp.Message += ", SUCCESSFULLY logged"
		}

		if procErr != nil {
			resp.ProcessError = procErr.Error()
			resp.Message += ", FAILED processing"
		} else {
			resp.Items = items
			resp.Message += ", SUCCESSFUL processing"
		}

		writeJSON(w, http.StatusOK, resp)

	}
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

func processBusinessData(ctx context.Context, body []byte, clients *app.Clients, logID *gen.LogID) (resultID *gen.LogID, results []*gen.VariableValue, err error) {
	var reqParsed requestJSON
	if err := json.Unmarshal(body, &reqParsed); err != nil {
		return nil, nil, fmt.Errorf("invalid JSON: %w", err)
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
		return nil, nil, fmt.Errorf("business logic error: %w", err)
	}

	return resp.LogID, results, nil
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
