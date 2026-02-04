package ws

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type MatchHub struct {
	upgrader websocket.Upgrader
	mu       sync.RWMutex
	clients  map[*websocket.Conn]string
}

func NewMatchHub() *MatchHub {
	return &MatchHub{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		clients: make(map[*websocket.Conn]string),
	}
}

type matchMessage struct {
	Action  string `json:"action"`
	MatchID string `json:"matchId"`
}

func (h *MatchHub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	h.register(conn, "")
	defer h.unregister(conn)

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			return
		}

		var msg matchMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}

		if msg.Action == "subscribe" && msg.MatchID != "" {
			h.register(conn, msg.MatchID)
			_ = conn.WriteJSON(map[string]string{"status": "subscribed", "matchId": msg.MatchID})
		}
		if msg.Action == "unsubscribe" {
			h.register(conn, "")
			_ = conn.WriteJSON(map[string]string{"status": "unsubscribed"})
		}
	}
}

func (h *MatchHub) Broadcast(matchID string, payload any) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for conn, subscribed := range h.clients {
		if subscribed != matchID {
			continue
		}
		_ = conn.WriteJSON(payload)
	}
}

func (h *MatchHub) register(conn *websocket.Conn, matchID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[conn] = matchID
}

func (h *MatchHub) unregister(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, conn)
	_ = conn.Close()
}
