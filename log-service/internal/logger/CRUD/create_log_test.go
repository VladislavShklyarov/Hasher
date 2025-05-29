package CRUD

import (
	"bytes"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log-service/gen"
	"strings"
	"testing"
	"time"
)

// bufWriteSyncer нужен для перехвата вывода logger
type bufWriteSyncer struct {
	*bytes.Buffer
}

func (bws bufWriteSyncer) Sync() error { return nil }

func TestWriteLogToFile(t *testing.T) {
	tests := []struct {
		name           string
		level          string
		serviceName    string
		expectInChan   bool
		expectChanFull bool
	}{
		{"debug level", "debug", "TestService", false, false},
		{"info level", "info", "TestService", false, false},
		{"warn level", "warn", "TestService", false, false},
		{"error level", "error", "TestService", false, false},
		{"default level", "unknown", "TestService", false, false},
		{"business-server ok", "info", "business-server", true, false},
		{"business-server full", "info", "business-server", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := bufWriteSyncer{&buf}

			encoderCfg := zap.NewProductionEncoderConfig()
			encoder := zapcore.NewJSONEncoder(encoderCfg)
			core := zapcore.NewCore(encoder, writer, zapcore.DebugLevel)
			logger := zap.New(core)

			logChan := make(chan *gen.LogEntry, 1)
			if tt.expectChanFull {
				logChan <- &gen.LogEntry{}
			}

			entry := &gen.LogEntry{
				ServiceName:   tt.serviceName,
				TimestampSend: time.Now().Add(-50 * time.Millisecond).UnixMilli(),
				Level:         tt.level,
				Message: &gen.StructuredMessage{
					Method: "GET",
					Path:   "/test/path",
					Body: []*gen.Operation{
						{Type: "calc", Op: "+", Var: "x", Left: "1", Right: "2"},
					},
					Result: &gen.OperationResponse{
						LogID: &gen.LogID{Id: "test-id"},
						Items: []*gen.VariableValue{{Var: "x", Value: 3}},
					},
				},
			}

			WriteLogToFile(logger, tt.level, "test-id", entry, logChan)

			logOutput := buf.String()

			if !strings.Contains(logOutput, `"id":"test-id"`) {
				t.Errorf("log output missing id: %s", logOutput)
			}
			if !strings.Contains(logOutput, `"source":"`+tt.serviceName+`"`) {
				t.Errorf("log output missing source: %s", logOutput)
			}
			if !strings.Contains(logOutput, "deliveryDelayMs") {
				t.Errorf("log output missing deliveryDelayMs: %s", logOutput)
			}
			if !strings.Contains(logOutput, "/test/path") {
				t.Errorf("log output missing path: %s", logOutput)
			}

			initialLen := 0
			if tt.expectChanFull {
				initialLen = 1
			}

			finalLen := len(logChan)

			if tt.expectInChan && finalLen != initialLen+1 {
				t.Error("expected log to be added to channel, but it wasn't")
			}
			if !tt.expectInChan && finalLen != initialLen {
				t.Error("log was unexpectedly added to channel")
			}

		})
	}
}
