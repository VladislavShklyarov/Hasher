package server

import (
	"os"
	"testing"
)

func TestNewLogManager(t *testing.T) {
	tests := []struct {
		name          string
		expectLoggers []string
		expectDir     string
	}{
		{
			name:          "default loggers and directory created",
			expectLoggers: []string{"HTTP-service", "business-service", "undefined-service"},
			expectDir:     "../log_files",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lm := NewLogManager()
			if lm == nil {
				t.Fatal("expected LogManager, got nil")
			}

			for _, key := range tt.expectLoggers {
				if _, ok := lm.Loggers[key]; !ok {
					t.Errorf("expected logger %q to be present", key)
				}
			}

			// Проверяем, что директория существует
			info, err := os.Stat(tt.expectDir)
			if err != nil {
				t.Errorf("expected log directory %q to exist, got error: %v", tt.expectDir, err)
			} else if !info.IsDir() {
				t.Errorf("expected %q to be a directory", tt.expectDir)
			}

			// Проверяем, что файлы для логов существуют
			for _, key := range tt.expectLoggers {
				filename := ""
				switch key {
				case "HTTP-service":
					filename = tt.expectDir + "/http_logs.json"
				case "business-service":
					filename = tt.expectDir + "/business_logs.json"
				case "undefined-service":
					filename = tt.expectDir + "/undefined_logs.json"
				}
				_, err := os.Stat(filename)
				if err != nil {
					t.Errorf("expected log file %q to exist, got error: %v", filename, err)
				}
			}
		})
	}
}
