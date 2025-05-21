package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func ValidateHttpRequest(r *http.Request) ([]byte, error) {
	if r.Method != http.MethodPost {
		return nil, errors.New("only POST requests allowed")
	}

	body, err := generalValidation(r)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func generalValidation(r *http.Request) ([]byte, error) {
	if r.ContentLength == 0 {
		return nil, errors.New("request body is empty")
	}

	if r.Header.Get("Content-Type") != "application/json" {
		return nil, errors.New("Content-Type must be application/json")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	if len(body) == 0 {
		return nil, errors.New("request body is empty")
	}

	var tmp interface{}
	if err := json.Unmarshal(body, &tmp); err != nil {
		return nil, errors.New("invalid JSON")
	}

	r.Body = io.NopCloser(bytes.NewBuffer(body))

	return body, nil
}
