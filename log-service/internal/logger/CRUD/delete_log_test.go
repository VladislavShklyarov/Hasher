package CRUD

import (
	"context"
	"io/fs"
	"log-service/gen"
	"os"
	"strings"
	"testing"
)

func TestDeleteLog(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		filePerm  fs.FileMode
		content   string
		logID     string
		expectMsg string
		expectOk  bool
		setupFile bool
	}{
		{
			name:      "log found and deleted",
			filename:  "delete_found.log",
			filePerm:  os.ModePerm,
			content:   `{"id":"abc123","msg":"hello world"}` + "\n" + `{"id":"def456","msg":"another log"}`,
			logID:     "abc123",
			expectMsg: "Log with id abc123 successfully deleted from delete_found.log",
			expectOk:  true,
			setupFile: true,
		},
		{
			name:      "log not found",
			filename:  "delete_notfound.log",
			filePerm:  os.ModePerm,
			content:   `{"id":"xyz999","msg":"other log"}`,
			logID:     "abc123",
			expectMsg: "log with id abc123 not found",
			expectOk:  false,
			setupFile: true,
		},
		{
			name:      "file does not exist",
			filename:  "missing.log",
			filePerm:  os.ModePerm,
			logID:     "abc123",
			expectMsg: "failed to open file:",
			expectOk:  false,
			setupFile: false,
		},
		{
			name:      "invalid json line but log found and deleted",
			filename:  "delete_invalid_json.log",
			filePerm:  os.ModePerm,
			content:   `{"id":"abc123","msg":}` + "\n" + `{"id":"def456","msg":"valid"}`,
			logID:     "def456",
			expectMsg: "Log with id def456 successfully deleted from delete_invalid_json.log",
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

			resp, err := lm.DeleteLog(context.Background(), req)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if resp.Success != tt.expectOk {
				t.Errorf("expected success: %v, got: %v", tt.expectOk, resp.Success)
			}

			if !strings.Contains(resp.Message, tt.expectMsg) {
				t.Errorf("expected message to contain %q, got: %q", tt.expectMsg, resp.Message)
			}
		})
	}
}

func TestUpdateLines(t *testing.T) {
	tests := []struct {
		name          string
		lines         []string
		idToDelete    string
		expectedLines []string
		expectErr     bool
	}{
		{
			name: "delete existing log line",
			lines: []string{
				`{"id":"abc123","msg":"hello"}`,
				`{"id":"def456","msg":"world"}`,
			},
			idToDelete: "abc123",
			expectedLines: []string{
				`{"id":"def456","msg":"world"}`,
			},
			expectErr: false,
		},
		{
			name: "log id not found",
			lines: []string{
				`{"id":"abc123","msg":"hello"}`,
			},
			idToDelete:    "nonexistent",
			expectedLines: nil,
			expectErr:     true,
		},
		{
			name: "invalid json line ignored",
			lines: []string{
				`{"id":"abc123","msg":}`, // invalid
				`{"id":"def456","msg":"valid"}`,
			},
			idToDelete: "def456",
			expectedLines: []string{
				`{"id":"abc123","msg":}`,
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updated, err := updateLines(tt.lines, tt.idToDelete)
			if (err != nil) != tt.expectErr {
				t.Fatalf("expected error: %v, got %v", tt.expectErr, err)
			}
			if !tt.expectErr {
				if len(updated) != len(tt.expectedLines) {
					t.Fatalf("expected %d lines, got %d", len(tt.expectedLines), len(updated))
				}
				for i := range updated {
					if updated[i] != tt.expectedLines[i] {
						t.Errorf("line %d: expected %q, got %q", i, tt.expectedLines[i], updated[i])
					}
				}
			}
		})
	}
}
