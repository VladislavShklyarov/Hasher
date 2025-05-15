package CRUD

import (
	"bytes"
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log-service/gen/logger"
	"strings"
	"testing"
	"time"
)

func TestLogEntry(t *testing.T) {
	tests := []struct {
		name  string
		level string
		id    string
		entry *gen.LogEntry
	}{
		{
			name:  "log level debug",
			level: "debug",
			id:    "test123",
			entry: &gen.LogEntry{
				Message:       "Debug message",
				ServiceName:   "TestService",
				TimestampSend: time.Now().Add(-50 * time.Millisecond).UnixMilli(),
			},
		}, {
			name:  "log level info",
			level: "info",
			id:    "test123",
			entry: &gen.LogEntry{
				Message:       "info message",
				ServiceName:   "TestService",
				TimestampSend: time.Now().Add(-50 * time.Millisecond).UnixMilli(),
			},
		}, {
			name:  "log level warn",
			level: "warn",
			id:    "test123",
			entry: &gen.LogEntry{
				Message:       "warn message",
				ServiceName:   "TestService",
				TimestampSend: time.Now().Add(-50 * time.Millisecond).UnixMilli(),
			},
		}, {
			name:  "log level error",
			level: "error",
			id:    "test123",
			entry: &gen.LogEntry{
				Message:       "error message",
				ServiceName:   "TestService",
				TimestampSend: time.Now().Add(-50 * time.Millisecond).UnixMilli(),
			},
		}, {
			name:  "log level default",
			level: "default",
			id:    "test123",
			entry: &gen.LogEntry{
				Message:       "default message",
				ServiceName:   "TestService",
				TimestampSend: time.Now().Add(-50 * time.Millisecond).UnixMilli(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			bws := bufWriteSyncer{&buf}

			encoderCfg := zap.NewProductionEncoderConfig()
			encoder := zapcore.NewJSONEncoder(encoderCfg)

			core := zapcore.NewCore(encoder, bws, zapcore.DebugLevel)
			logger := zap.New(core)

			WriteLogToFile(logger, tt.level, tt.id, tt.entry)

			logs := buf.String()

			if !strings.Contains(logs, tt.entry.Message) {
				t.Errorf("expected message '%s' in logs, got: %s", tt.entry.Message, logs)
			}
			if !strings.Contains(logs, tt.id) {
				t.Errorf("expected id '%s' in logs, got: %s", tt.id, logs)
			}
			if !strings.Contains(logs, "deliveryDelayMs") {
				t.Errorf("expected delivery delay field in logs, got: %s", logs)
			}
		})
	}
}

func TestHandleIncomingLog(t *testing.T) {
	type testCase struct {
		name         string
		serviceName  string
		expectLogger string
		expectErr    bool
		ctxCancelled bool
	}

	tests := []testCase{
		{
			name:         "Known service",
			serviceName:  "HTTP-service",
			expectLogger: "HTTP-service",
		},
		{
			name:         "Unknown service uses fallback",
			serviceName:  "UnknownService",
			expectLogger: "undefined-service",
		},
		{
			name:         "Context cancelled",
			serviceName:  "HTTP-service",
			expectErr:    true,
			ctxCancelled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bufs := make(map[string]*bytes.Buffer)
			loggers := make(map[string]*zap.Logger)

			for _, svc := range []string{"HTTP-service", "undefined-service"} {
				buf := &bytes.Buffer{}
				bufs[svc] = buf

				ws := zapcore.AddSync(buf)
				encoderCfg := zap.NewProductionEncoderConfig()
				encoder := zapcore.NewJSONEncoder(encoderCfg)
				core := zapcore.NewCore(encoder, ws, zapcore.DebugLevel)
				loggers[svc] = zap.New(core)
			}

			lm := LogManager{
				Loggers: loggers,
			}

			entry := &gen.LogEntry{
				Message:       "test log message",
				ServiceName:   tt.serviceName,
				TimestampSend: time.Now().Add(-50 * time.Millisecond).UnixMilli(),
				Level:         "info",
			}

			ctx := context.Background()
			if tt.ctxCancelled {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel()
			}

			resp, err := lm.HandleIncomingLog(ctx, entry)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			buf := bufs[tt.expectLogger]
			logOutput := buf.String()
			if !strings.Contains(logOutput, entry.Message) {
				t.Errorf("log output missing message: got %s", logOutput)
			}
			if resp.GetId() == nil || resp.GetId().GetId() == "" {
				t.Errorf("expected non-empty response ID")
			}
		})
	}
}

type bufWriteSyncer struct {
	*bytes.Buffer
}

func (bws bufWriteSyncer) Sync() error {
	// Обычно ничего не делаем при Sync для буфера
	return nil
}
