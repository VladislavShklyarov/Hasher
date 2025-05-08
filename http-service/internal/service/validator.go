package service

import (
	"encoding/json"
	"io"
	"net/http"
)

func ValidateHttpRequest(w http.ResponseWriter, r *http.Request) (body []byte) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests allowed", http.StatusMethodNotAllowed)
		return nil
	}

	if r.ContentLength == 0 {
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return nil
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return nil
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if len(body) == 0 {
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return nil
	}

	if err != nil {
		http.Error(w, "Failed to read request body. Check whether body is correct", http.StatusInternalServerError)
		return nil
	}

	var tmp interface{}
	if err := json.Unmarshal(body, &tmp); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return nil
	}

	return body
}
