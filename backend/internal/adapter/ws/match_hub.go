package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type MatchHub struct {
	upgrader websocket.Upgrader
	mu       sync.RWMutex
	clients  map[*websocket.Conn]string
}

// AllowedWebSocketOrigins defines the list of trusted origins for WebSocket connections
var AllowedWebSocketOrigins = []string{
	"http://localhost:3000",
	"http://localhost:5173",
	"http://127.0.0.1:3000",
	"http://127.0.0.1:5173",
	// Add production domain(s) here
}

func isWSOriginAllowed(origin string) bool {
	for _, allowed := range AllowedWebSocketOrigins {
		if origin == allowed {
			return true
		}
	}
	return false
}

func NewMatchHub() *MatchHub {
	return &MatchHub{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				if origin == "" {
					return true // Allow requests without Origin header (same-origin)
				}
				allowed := isWSOriginAllowed(origin)
				if !allowed {
					log.Printf("WebSocket origin rejected: %s", origin)
				}
				return allowed
			},
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
