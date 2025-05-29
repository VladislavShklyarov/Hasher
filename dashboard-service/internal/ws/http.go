package wscd

import (
	"dashboard-service/internal/config"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

// Мапа для активных клиентов

type Clients struct {
	Clients   map[*websocket.Conn]bool
	clientsMu sync.Mutex
}

func (c *Clients) AddClient(conn *websocket.Conn) {
	c.clientsMu.Lock()
	c.Clients[conn] = true
	c.clientsMu.Unlock()
}

func (c *Clients) DeleteClient(conn *websocket.Conn) {
	c.clientsMu.Lock()
	delete(c.Clients, conn)
	c.clientsMu.Unlock()
}

func StartWebSocket(clients *Clients, cfg *config.Config) {

	http.HandleFunc("/ws", handleConnections(clients)) // обработчик WebSocket

	http.HandleFunc("/", htmlHandler()) // обработчик главной страницы

	fmt.Printf("web socket is running on %s\n", cfg.WsAddr)
	err := http.ListenAndServe(cfg.WsAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func htmlHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>WebSocket Image and Log Viewer</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            padding: 20px;
            text-align: center;
        }
        img {
            max-width: 90%;
            height: auto;
            border: 1px solid #ccc;
            margin-bottom: 20px;
        }
        pre {
            text-align: left;
            max-width: 90%;
            margin: 0 auto;
            background: #f4f4f4;
            padding: 10px;
            border-radius: 5px;
            overflow-x: auto;
            white-space: pre-wrap;
        }
    </style>
</head>
<body>
    <h1>WebSocket Image and Log Viewer</h1>
    <img id="wsImage" alt="Waiting for image..." />
    <pre id="wsText">Waiting for logs...</pre>

    <script>
        const ws = new WebSocket("ws://" + location.host + "/ws");
        ws.binaryType = "arraybuffer";

        ws.onmessage = function(event) {
            if (typeof event.data === "string") {
                // Text message (e.g., protobuf String or JSON)
                try {
                    const json = JSON.parse(event.data);
                    document.getElementById("wsText").textContent = JSON.stringify(json, null, 2);
                } catch (e) {
                    document.getElementById("wsText").textContent = event.data;
                }
            } else {
                // Binary message (image)
                const blob = new Blob([event.data], { type: "image/jpeg" });
                const url = URL.createObjectURL(blob);
                document.getElementById("wsImage").src = url;
            }
        };

        ws.onclose = function() {
            console.log("WebSocket connection closed");
        };

        ws.onerror = function(err) {
            console.error("WebSocket error:", err);
        };
    </script>
</body>
</html>
	`
		w.Write([]byte(html))
	}
}

func BroadcastImage(imageData []byte, clients *Clients) {
	clients.clientsMu.Lock()
	defer clients.clientsMu.Unlock()

	for client := range clients.Clients {
		err := client.WriteMessage(websocket.BinaryMessage, imageData)
		if err != nil {
			log.Println("Broadcast error:", err)
			client.Close()
			delete(clients.Clients, client)
		}
	}
}

func BroadcastOperation(operationData []byte, clients *Clients) {
	clients.clientsMu.Lock()
	defer clients.clientsMu.Unlock()

	for client := range clients.Clients {
		err := client.WriteMessage(websocket.TextMessage, operationData)
		if err != nil {
			log.Println("Broadcast error:", err)
			client.Close()
			delete(clients.Clients, client)
		}
	}
}
