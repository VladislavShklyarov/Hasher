package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	gen "http-service/gen/logger"
	"http-service/internal/app"
	"net/http"
)

func ReadLogHandler(clients *app.Clients) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id := r.URL.Query().Get("id")
		filename := r.URL.Query().Get("filename")

		if id == "" || filename == "" {
			writeJSONError(w, false, http.StatusBadRequest, "Missing query parameters", "id and filename are required")
			return
		}

		readLogResponse, err := clients.LogClient.ReadLogGRPC(id, filename)
		if err != nil {
			writeJSONError(w, readLogResponse.GetSuccess(), http.StatusInternalServerError, "Failed to retrieve log due to server mailfunction", err.Error())
			return
		}

		if !readLogResponse.GetSuccess() {
			writeJSONError(w, readLogResponse.GetSuccess(), http.StatusInternalServerError, "Failed to retrive log", readLogResponse.GetError())
			return
		}

		readLogResponseJSON, err := parseReadResponse(readLogResponse)
		if err != nil {
			writeJSONError(w, readLogResponse.GetSuccess(), http.StatusInternalServerError, "Failed to marshal log:", err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(readLogResponseJSON)
	}

}

func parseReadResponse(readLogResponse *gen.LogReadingResponse) (responseBytes []byte, err error) {

	var logEntry LogEntry

	if readLogResponse.GetSuccess() && readLogResponse.Log != "" {
		err = json.Unmarshal([]byte(readLogResponse.Log), &logEntry)
		if err != nil {
			return nil, fmt.Errorf("failed to parse nested log entry: %w", err)
		}
	}

	response := ReadResponse{
		Success: readLogResponse.GetSuccess(),
		Log:     logEntry,
	}

	if !readLogResponse.GetSuccess() {
		response.Error = readLogResponse.GetError()
	}

	responseBytes, err = json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal final response: %w", err)
	}
	fmt.Println(response)
	return responseBytes, nil

}

type ReadResponse struct {
	Success bool     `json:"success"`
	Log     LogEntry `json:"log,omitempty"`
	Error   string   `json:"error,omitempty"`
}
