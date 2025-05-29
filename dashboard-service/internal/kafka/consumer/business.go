package consumer

import (
	wscd "dashboard-service/internal/ws"
	"encoding/base64"
	"log"
)

type BizHandler struct{}

func (b *BizHandler) Handle(encoded []byte) []byte {
	data, err := decodeBase64(encoded)
	if err != nil {
		log.Printf("Ошибка декодирования base64: %v\n", err)
		return nil
	}
	return data
}

func (b *BizHandler) Broadcast(data []byte, clients *wscd.Clients) {
	wscd.BroadcastImage(data, clients)
}

func decodeBase64(encoded []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(encoded))
}
