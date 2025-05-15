package utils

import (
	"strings"
	"testing"
)

func TestGenerateID(t *testing.T) {
	tests := []struct {
		name      string
		length    int
		expectLen int
	}{
		{
			name:      "generate id of length 10",
			length:    10,
			expectLen: 10,
		},
		{
			name:      "generate id of length 0",
			length:    0,
			expectLen: 0,
		},
		{
			name:      "generate id of length 100",
			length:    100,
			expectLen: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := GenerateID(tt.length)
			if len(id) != tt.expectLen {
				t.Errorf("expected length %d, got %d", tt.expectLen, len(id))
			}

			// Проверка допустимых символов
			const allowed = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
			for _, r := range id {
				if !strings.ContainsRune(allowed, r) {
					t.Errorf("id contains invalid character: %q", r)
				}
			}
		})
	}
}
