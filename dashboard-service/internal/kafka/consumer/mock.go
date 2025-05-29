package consumer

import wscd "dashboard-service/internal/ws"

type MockHandler struct {
	HandleFunc    func(msg []byte) []byte
	BroadcastFunc func(data []byte, clients *wscd.Clients)
}

func (m *MockHandler) Handle(msg []byte) []byte {
	if m.HandleFunc != nil {
		return m.HandleFunc(msg)
	}
	return nil
}

func (m *MockHandler) Broadcast(data []byte, clients *wscd.Clients) {
	if m.BroadcastFunc != nil {
		m.BroadcastFunc(data, clients)
	}
}
