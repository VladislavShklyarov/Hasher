package CRUD

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log-service/gen"
	"os"
	"path/filepath"
)

func (*LogManager) ReadLog(ctx context.Context, logInfo *gen.LogInfo) (*gen.LogReadingResponse, error) {

	filename := logInfo.GetFilename()

	file, err := OpenFile(filename)
	if err != nil {
		return writeWrongReadResponse("failed to open file", err), nil
	}
	defer file.Close()

	log, err := FindLog(file, logInfo.GetId())
	if err != nil {
		return writeWrongReadResponse(fmt.Sprintf("failed to find log in %s", filename), err), nil
	}

	return &gen.LogReadingResponse{
		Success: true,
		Log:     log,
		Error:   "",
	}, nil

}

func OpenFile(filename string) (file *os.File, err error) {
	filePath := filepath.Join("../log_files/", filename)
	file, err = os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func FindLog(file *os.File, id string) (log string, err error) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		var logId struct {
			ID string `json:"id"`
		}

		if err := json.Unmarshal([]byte(line), &logId); err != nil {
			continue // Пропускаем невалидные JSON строки
		}

		if logId.ID == id {
			return line, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error while scanning file: %w", err)
	}

	return "", fmt.Errorf("log with id %s not found", id)
}
func writeWrongReadResponse(msg string, err error) *gen.LogReadingResponse {
	var errMsg string
	if err != nil {
		errMsg = fmt.Sprintf("%s: %s", msg, err.Error())
	} else {
		errMsg = msg
	}

	return &gen.LogReadingResponse{
		Success: false,
		Error:   errMsg,
	}
}
