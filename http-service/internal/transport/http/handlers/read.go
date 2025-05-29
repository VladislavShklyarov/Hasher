package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	gen "http-service/gen"
	"http-service/internal/app"
	"net/http"
)

func ReadLogHandler(clients *app.Clients) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id := r.URL.Query().Get("id")
		filename := r.URL.Query().Get("filename")

		if id == "" || filename == "" {
			writeJSON(w, http.StatusBadRequest, "Missing query parameters: id and filename are required")
			return
		}

		readLogResponse, err := clients.LogClient.ReadLogGRPC(id, filename)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, "Failed to retrieve log due to server mailfunction: "+err.Error())
			return
		}

		if !readLogResponse.GetSuccess() {
			writeJSON(w, http.StatusInternalServerError, "Failed to retrive log: "+readLogResponse.GetError())
			return
		}

		readLogResponseJSON, err := parseReadResponse(readLogResponse)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, "Failed to marshal log: "+err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(readLogResponseJSON)
	}

}

func parseReadResponse(readLogResponse *gen.LogReadingResponse) ([]byte, error) {
	if !readLogResponse.Success {
		resp := ReadResponse{
			Success: false,
			Log:     LogEntry{},
			Error:   readLogResponse.Error,
		}
		return json.Marshal(resp)
	}

	var logEntry LogEntry

	err := json.Unmarshal([]byte(readLogResponse.Log), &logEntry)
	if err != nil {
		return nil, fmt.Errorf("failed to parse nested log entry: %w", err)
	}

	var parsedMessage interface{}
	err = json.Unmarshal(logEntry.Message, &parsedMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to parse nested message JSON: %w", err)
	}

	type LogResponse struct {
		Level             string      `json:"level"`
		Msg               string      `json:"msg"`
		Id                string      `json:"id"`
		Message           interface{} `json:"message"`
		Source            string      `json:"source"`
		TimestampSend     int64       `json:"timestamp_send"`
		TimestampReceived int64       `json:"timestamp_received"`
		DeliveryDelayMs   string      `json:"deliveryDelayMs"`
	}

	response := struct {
		Success bool        `json:"success"`
		Log     LogResponse `json:"log,omitempty"`
		Error   string      `json:"error,omitempty"`
	}{
		Success: true,
		Log: LogResponse{
			Level:             logEntry.Level,
			Msg:               logEntry.Msg,
			Id:                logEntry.ID,
			Message:           parsedMessage,
			Source:            logEntry.Source,
			TimestampSend:     logEntry.TimestampSend,
			TimestampReceived: logEntry.TimestampRecv,
			DeliveryDelayMs:   logEntry.DeliveryDelayMs,
		},
	}

	return json.Marshal(response)
}

type ReadResponse struct {
	Success bool     `json:"success"`
	Log     LogEntry `json:"log,omitempty"`
	Error   string   `json:"error,omitempty"`
}
