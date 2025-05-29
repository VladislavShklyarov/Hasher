package wscd

import (
	"log"
	"net/http"
)

func handleConnections(clients *Clients) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}
		defer conn.Close()
		clients.AddClient(conn)
		log.Println("New WebSocket clients connected:")

		for {
			_, _, err = conn.ReadMessage()
			if err != nil {
				log.Println("Client disconnected: ", err)
				clients.DeleteClient(conn)
			}
		}
	}
}
