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
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	h.register(conn, "")
	defer h.unregister(conn)

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			return
		}

		var msg matchMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Printf("WebSocket message unmarshal error: %v", err)
			// Send error frame to client
			if writeErr := conn.WriteJSON(map[string]string{"error": "invalid message format"}); writeErr != nil {
				log.Printf("WebSocket write error (error response): %v", writeErr)
			}
			continue
		}

		if msg.Action == "subscribe" && msg.MatchID != "" {
			h.register(conn, msg.MatchID)
			if err := conn.WriteJSON(map[string]string{"status": "subscribed", "matchId": msg.MatchID}); err != nil {
				log.Printf("WebSocket write error (subscribe response): %v", err)
			}
		} else if msg.Action == "unsubscribe" {
			h.register(conn, "")
			if err := conn.WriteJSON(map[string]string{"status": "unsubscribed"}); err != nil {
				log.Printf("WebSocket write error (unsubscribe response): %v", err)
			}
		} else {
			log.Printf("WebSocket unknown action: %s", msg.Action)
			if err := conn.WriteJSON(map[string]string{"error": "unknown action"}); err != nil {
				log.Printf("WebSocket write error (unknown action response): %v", err)
			}
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
		if err := conn.WriteJSON(payload); err != nil {
			log.Printf("WebSocket broadcast error (matchID: %s): %v", matchID, err)
		}
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
