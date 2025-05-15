package log

import (
	"http-service/gen/logger"
	"net/http"
)

type LogClient struct {
	LoggerClient gen.LoggerClient
}
type LogClientInterface interface {
	ReadLogGRPC(id, filename string) (*gen.LogReadingResponse, error)
	DeleteLogGRPC(id, filename string) (*gen.LogDeletionResponse, error)
	LogDataGRPC(r *http.Request) (id string, err error)
}

var _ LogClientInterface = (*LogClient)(nil) // compile-time check
