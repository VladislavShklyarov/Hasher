package main

import (
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"http-service/internal/app"
	grpcBiz "http-service/internal/client/grpc/business"
	grpcLog "http-service/internal/client/grpc/log"
	"http-service/internal/config"
	"http-service/internal/server"
	"http-service/internal/signals"
	_ "http-service/internal/transport/http"
	"net/http"
)

// @title			Swagger Example API
// @version		1.0
// @description	This is a sample server celler server
// @termOfService	http://swagger.io/terms
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Load()

	clients := &app.Clients{
		LogClient:      grpcLog.CreateLogClient(cfg),
		BusinessClient: grpcBiz.CreateBusinessClient(cfg),
	}

	go server.RunHttpServer(clients, cfg)

	signals.WaitForShutdown(ctx, cancel)
}

// ProcessDataHandler обрабатывает POST запрос для отправки данных и их логирования.
//
//	@Summary		Process and log data
//	@Description	Receives data from the request body, validates it, logs it via gRPC,
//
// and returns a success or failure message based on the result of the logging process.
//
//	@Tags			data
//	@Accept			json
//	@Produce		text/plain
//	@Param			body	body		string	true	"Data to be processed"
//	@Success		200		{string}	string	"Log written successfully"
//	@Failure		400		{string}	string	"Invalid request"
//	@Failure		500		{string}	string	"Internal server error"
//	@Router			/process [post]

// ReadLogHandler retrieves a log entry by ID and filename.
//
//	@Summary		Retrieve log entry
//	@Description	Retrieves a log entry based on its ID and filename passed as query parameters.
//	@Tags			logs
//	@Accept			json
//	@Produce		json
//	@Param			id			query		string	true	"Log entry ID"
//	@Param			filename	query		string	true	"Filename containing the log"
//	@Success		200			{object}	LogEntry
//	@Failure		400			{object}	ErrorResponse	"Missing query parameters"
//	@Failure		500			{object}	ErrorResponse	"Failed to retrieve or marshal log"
//	@Router			/getLog [get]
func ReadLogHandler(LogClient *grpcLog.LogClient) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id := r.URL.Query().Get("id")
		filename := r.URL.Query().Get("filename")

		if id == "" || filename == "" {
			http.Error(w, "Missing query parameters", http.StatusBadRequest)
			return
		}

		logEntry, err := LogClient.ReadLogGRPC(id, filename)
		if err != nil {
			http.Error(w, "Failed to retrieve log: "+err.Error(), http.StatusInternalServerError)
			return
		}

		logEntryJSON, err := json.Marshal(logEntry)
		if err != nil {
			http.Error(w, "Failed to marshal log entry: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(logEntryJSON)
	}

}

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
	Level     string `json:"level"`
}

// ErrorResponse represents a standard error response.
// swagger:model
type ErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
	Reason string `json:"reason"`
}

// DeleteLogHandler handles the deletion of a log entry by its ID and filename.
//
//	@Summary		Delete a log entry
//	@Description	Deletes a log entry identified by ID from the specified filename.
//	@Tags			logs
//	@Accept			json
//	@Produce		json
//	@Param			id			query		string	true	"Log ID"
//	@Param			filename	query		string	true	"Filename where the log is stored"
//	@Success		200			{object}	DeleteResponse
//	@Failure		400			{object}	ErrorResponse	"Missing query parameters"
//	@Failure		500			{object}	ErrorResponse	"Internal server error"
//	@Router			/deleteLog [delete]
func DeleteLogHandler(LogClient *grpcLog.LogClient) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id := r.URL.Query().Get("id")
		filename := r.URL.Query().Get("filename")

		if id == "" || filename == "" {
			writeJSONError(w, http.StatusBadRequest, "Missing query parameters", "id and filename are required")
			return
		}

		deleteResponse, err := LogClient.DeleteLogGRPC(id, filename)

		if err != nil || deleteResponse == nil {
			errMsg := "Unknown error"
			if err != nil {
				errMsg = err.Error()
			}
			response := DeleteResponse{
				Success: false,
				Error:   errMsg,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		response := DeleteResponse{
			Success: deleteResponse.Success,
			Message: deleteResponse.Message,
			Error:   "",
		}

		if !deleteResponse.Success {
			if deleteResponse.Message == "" {
				response.Message = "Deletion failed"
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
func writeJSONError(w http.ResponseWriter, statusCode int, errorName string, reason string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Status: statusCode, Error: errorName, Reason: reason})
}

// DeleteResponse represents the result of a delete operation.
// Used in successful responses.
// swagger:model
type DeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error"`
}
