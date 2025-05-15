package handlers

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"http-service/internal/app"
	"http-service/internal/utils"
	"net/http"
)

func ProcessDataHandler(clients *app.Clients) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		_, err := utils.ValidateHttpRequest(r)
		defer r.Body.Close()
		if err != nil {
			writeJSONError(w, false, http.StatusBadRequest, "Invalid request body", err.Error())
			return
		}

		logId, logErr := clients.LogClient.LogDataGRPC(r)

		message := "Request received successfully"
		if logErr != nil {
			message += ", but log server is temporarily unavailable"
		} else {
			message += ". Logged successfully."
		}

		writeJSONSuccessResponce(w, true, http.StatusOK, message, logId)
	}
}

type SuccesResponse struct {
	Success bool   `json:"success"`
	Status  int    `json:"status"`
	Message string `json:"message"`
	Id      string `json:"id"`
}

func writeJSONSuccessResponce(w http.ResponseWriter, success bool, statusCode int, message string, id string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(SuccesResponse{Success: success, Status: statusCode, Message: message, Id: id})
}
