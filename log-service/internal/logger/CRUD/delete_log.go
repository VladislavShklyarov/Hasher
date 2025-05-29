package CRUD

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log-service/gen"
	"path/filepath"

	"os"
	"strings"
)

func (*LogManager) DeleteLog(ctx context.Context, logInfo *gen.LogInfo) (*gen.LogDeletionResponse, error) {
	filename := logInfo.GetFilename()
	id := logInfo.GetId()

	filePath := filepath.Join("../log_files/", filename)

	file, err := OpenFile(filename)
	if err != nil {
		return writeWrongDeleteResponse("failed to open file: " + err.Error()), nil
	}
	defer file.Close()

	lines, err := readLines(file)
	if err != nil {
		return writeWrongDeleteResponse("failed to scan file: " + ": " + err.Error()), nil
	}

	updatedLines, err := updateLines(lines, id)
	if err != nil {
		return writeWrongDeleteResponse(err.Error()), nil
	}

	err = rewriteFile(filePath, updatedLines)
	if err != nil {
		return writeWrongDeleteResponse("failed to overwrite file: " + err.Error()), nil
	}

	return &gen.LogDeletionResponse{
		Success: true,
		Message: "Log with id " + id + " successfully deleted from " + filename}, nil

}

func writeWrongDeleteResponse(err string) *gen.LogDeletionResponse {
	return &gen.LogDeletionResponse{
		Success: false,
		Message: err,
	}
}

func readLines(file *os.File) ([]string, error) {

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func updateLines(lines []string, id string) (updatedLines []string, err error) {
	var found bool

	for _, line := range lines {
		var entry map[string]any
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			updatedLines = append(updatedLines, line)
			continue
		}
		if entryID, ok := entry["id"].(string); ok && entryID == id {
			found = true
			continue
		}
		updatedLines = append(updatedLines, line)
	}
	if !found {
		err = fmt.Errorf("log with id %s not found", id)
		return nil, err
	}
	return updatedLines, err
}

func rewriteFile(filePath string, updatedLines []string) (err error) {
	err = os.WriteFile(filePath, []byte(strings.Join(updatedLines, "\n")+"\n"), 0644)
	if err != nil {
		return err
	}
	return nil
}
