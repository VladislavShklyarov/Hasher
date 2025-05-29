package consumer

import (
	"dashboard-service/gen"
	wscd "dashboard-service/internal/ws"
	"fmt"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type LogHandler struct{}

func (l *LogHandler) Handle(data []byte) []byte {
	return data
}

func (l *LogHandler) Broadcast(data []byte, clients *wscd.Clients) {
	var message gen.StructuredMessage
	err := proto.Unmarshal(data, &message)
	if err != nil {
		fmt.Println("Failed to parse: ", err)
	}

	jsonData, err := protojson.Marshal(&message)
	if err != nil {
		fmt.Println("Failed to convert to JSON:", err)
		return
	}

	wscd.BroadcastOperation(jsonData, clients)

}
