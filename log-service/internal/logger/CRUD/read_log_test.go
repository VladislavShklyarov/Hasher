package CRUD

import (
	"context"
	"errors"
	"io/fs"
	"log-service/gen/logger"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFindLog(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		id        string
		expect    string
		expectErr bool
	}{
		{name: "log found",
			content: `{"id":"abc123","msg":"test message"}
{"id":"xyz789","msg":"another message"}`,
			id:        "abc123",
			expect:    `{"id":"abc123","msg":"test message"}`,
			expectErr: false,
		},
		{name: "invalid json line found",
			content: `{"id":"abc123","msg":}
{"id":"xyz789","msg":"another message"}`,
			id:        "abc123",
			expect:    ``,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := createTempFile(t, tt.content)
			defer file.Close()

			got, err := findLog(file, tt.id)
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
			if strings.TrimSpace(got) != strings.TrimSpace(tt.expect) {
				t.Errorf("expected log line: %s, got: %s", tt.expect, got)
			}
		})
	}

}

func TestWriteWrongReadResponse(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		err      error
		expected *gen.LogReadingResponse
	}{
		{
			name:    "simple error",
			message: "file error",
			err:     errors.New("file not found"),
			expected: &gen.LogReadingResponse{
				Success: false,
				Error:   "file error: file not found",
			},
		},
		{
			name:    "empty message",
			message: "",
			err:     errors.New("permission denied"),
			expected: &gen.LogReadingResponse{
				Success: false,
				Error:   ": permission denied",
			},
		},
		{
			name:    "nil error",
			message: "validation failed",
			err:     nil,
			expected: &gen.LogReadingResponse{
				Success: false,
				Error:   "validation failed", // ← исправлено
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := writeWrongReadResponse(tt.message, tt.err)

			if got.Success != tt.expected.Success {
				t.Errorf("Success got = %v, want %v", got.Success, tt.expected.Success)
			}

			if got.Error != tt.expected.Error {
				t.Errorf("Error got = %v, want %v", got.Error, tt.expected.Error)
			}
		})
	}
}

func TestOpenFile(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		filePerm  fs.FileMode
		content   string
		create    bool
		expectErr bool
	}{
		{
			name:      "file exists",
			filename:  "test_log_exists.log",
			filePerm:  os.ModePerm,
			content:   "test content",
			create:    true,
			expectErr: false,
		},
		{
			name:      "file does not exist",
			filename:  "missing_file.log",
			create:    false,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var filePath string
			if tt.create {
				filePath = createTempLogFile(t, tt.filePerm, tt.filename, tt.content)
				defer cleanupTempLogFile(t, filePath)
			}

			file, err := openFile(tt.filename)
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
			if file != nil {
				defer file.Close()
			}
		})
	}
}

func TestReadLog(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		filePerm  fs.FileMode
		content   string
		logID     string
		expectLog string
		expectErr bool
		expectOk  bool
		setupFile bool
	}{
		{
			name:      "log found in existing file",
			filename:  "readlog_found.log",
			filePerm:  os.ModePerm,
			content:   `{"id":"abc123","msg":"hello world"}` + "\n" + `{"id":"def456","msg":"another log"}`,
			logID:     "abc123",
			expectLog: `{"id":"abc123","msg":"hello world"}`,
			expectErr: false,
			expectOk:  true,
			setupFile: true,
		},
		{
			name:      "log not found",
			filename:  "readlog_notfound.log",
			filePerm:  os.ModePerm,
			content:   `{"id":"xyz999","msg":"other log"}`,
			logID:     "abc123",
			expectLog: "",
			expectErr: false,
			expectOk:  false,
			setupFile: true,
		},
		{
			name:      "file does not exist",
			filename:  "missing.log",
			filePerm:  os.ModePerm,
			logID:     "abc123",
			expectLog: "",
			expectErr: false,
			expectOk:  false,
			setupFile: false,
		},
		{ // Проверка на устойсивость к ошибке, если перед искомым логом невалидный json
			name:      "invalid json in file",
			filename:  "invalid_json.log",
			filePerm:  os.ModePerm,
			content:   `{"id":"abc123","msg":}` + "\n" + `{"id":"def456","msg":"valid"}`,
			logID:     "def456",
			expectLog: `{"id":"def456","msg":"valid"}`,
			expectErr: false,
			expectOk:  true,
			setupFile: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFile {
				path := createTempLogFile(t, tt.filePerm, tt.filename, tt.content)
				defer cleanupTempLogFile(t, path)
			}

			lm := LogManager{}
			req := &gen.LogInfo{
				Filename: tt.filename,
				Id:       tt.logID,
			}

			resp, err := lm.ReadLog(context.Background(), req)
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
			if resp.Success != tt.expectOk {
				t.Errorf("expected success: %v, got: %v", tt.expectOk, resp.Success)
			}
			if strings.TrimSpace(resp.Log) != strings.TrimSpace(tt.expectLog) {
				t.Errorf("expected log: %s, got: %s", tt.expectLog, resp.Log)
			}
		})
	}
}

func createTempFile(t *testing.T, content string) *os.File {
	tmpFile, err := os.CreateTemp("", "test_log_*.log")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	t.Cleanup(func() { os.Remove(tmpFile.Name()) })

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	if _, err := tmpFile.Seek(0, 0); err != nil {
		t.Fatalf("failed to seek to start of file: %v", err)
	}

	return tmpFile
}

func createTempLogFile(t *testing.T, perm fs.FileMode, filename, content string) string {
	t.Helper()

	dir := "../log_files"
	err := os.MkdirAll(dir, perm)
	if err != nil {
		t.Fatalf("failed to create log_files dir: %v", err)
	}

	path := filepath.Join(dir, filename)
	err = os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write temp log file: %v", err)
	}

	return path
}

func cleanupTempLogFile(t *testing.T, path string) {
	t.Helper()
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		t.Errorf("failed to remove temp log file: %v", err)
	}
}
