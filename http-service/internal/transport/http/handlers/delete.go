package handlers

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"http-service/internal/app"
	"net/http"
)

func DeleteLogHandler(clients *app.Clients) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id := r.URL.Query().Get("id")
		filename := r.URL.Query().Get("filename")

		if id == "" || filename == "" {
			writeJSONError(w, false, http.StatusBadRequest, "Missing query parameters", "id and filename are required")
			return
		}

		deleteResponse, err := clients.LogClient.DeleteLogGRPC(id, filename)

		if err != nil {
			writeJSONError(w, deleteResponse.GetSuccess(), http.StatusInternalServerError, "Failed to delete log due to server mailfunction", err.Error())
			return
		}

		if !deleteResponse.GetSuccess() {
			writeJSONError(w, deleteResponse.GetSuccess(), http.StatusInternalServerError, "Failed to delete log", deleteResponse.Message)
			return
		}

		deleteLogResponseJSON, err := json.Marshal(deleteResponse)
		if err != nil {
			writeJSONError(w, deleteResponse.GetSuccess(), http.StatusInternalServerError, "Failed to marshal log:", err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(deleteLogResponseJSON)

	}
}

type DeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
