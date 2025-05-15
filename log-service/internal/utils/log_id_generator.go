package utils

import "math/rand"

func GenerateID(logLength int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	logLengthBytes := make([]rune, logLength)
	for i := range logLengthBytes {
		logLengthBytes[i] = rune(letters[rand.Intn(len(letters))])
	}
	return string(logLengthBytes)
}
