package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"empoweredpixels/internal/mcp"
)

// MCPRESTHandler handles REST endpoints for MCP agents
type MCPRESTHandler struct {
	mcp        *mcp.MCPHandler
	auditLogger *mcp.AuditLogger
	filter     *mcp.FairnessFilter
}

// NewMCPRESTHandler creates a new MCP REST handler
func NewMCPRESTHandler(mcp *mcp.MCPHandler, auditLogger *mcp.AuditLogger, filter *mcp.FairnessFilter) *MCPRESTHandler {
	return &MCPRESTHandler{
		mcp:        mcp,
		auditLogger: auditLogger,
		filter:     filter,
	}
}

// getAPIKey extracts the API key from the request header
func (h *MCPRESTHandler) getAPIKey(r *http.Request) string {
	return r.Header.Get("X-EP-AI-KEY")
}

// authenticate checks the API key and returns the userID
func (h *MCPRESTHandler) authenticate(ctx context.Context, apiKey string) (int64, error) {
	return h.mcp.Authenticate(ctx, apiKey)
}

// GameState returns the current game state (lobbies and matches)
func (h *MCPRESTHandler) GameState(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	apiKey := h.getAPIKey(r)
	endpoint := "/mcp/game-state"

	// Apply rate limiting
	allowed, err := h.filter.Allow(apiKey)
	if !allowed {
		h.auditLogger.LogAction(ctx, "game_state_denied", apiKey, 0, endpoint, nil, false, err.Error())
		http.Error(w, err.Error(), http.StatusTooManyRequests)
		return
	}

	// Authenticate
	userID, err := h.authenticate(ctx, apiKey)
	if err != nil {
		h.auditLogger.LogAction(ctx, "game_state_auth_fail", apiKey, 0, endpoint, nil, false, err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Get game state
	state, err := h.mcp.GetGameState(ctx, userID)
	if err != nil {
		h.auditLogger.LogAction(ctx, "game_state_error", apiKey, userID, endpoint, nil, false, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log success
	h.auditLogger.LogAction(ctx, "game_state_success", apiKey, userID, endpoint, nil, true, "")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

// PlayerStats returns comprehensive stats for a fighter
func (h *MCPRESTHandler) PlayerStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	apiKey := h.getAPIKey(r)
	playerID := r.URL.Query().Get("id")
	endpoint := "/mcp/player/stats"

	// Apply rate limiting
	allowed, err := h.filter.Allow(apiKey)
	if !allowed {
		h.auditLogger.LogAction(ctx, "player_stats_denied", apiKey, 0, endpoint, nil, false, err.Error())
		http.Error(w, err.Error(), http.StatusTooManyRequests)
		return
	}

	// Authenticate
	userID, err := h.authenticate(ctx, apiKey)
	if err != nil {
		h.auditLogger.LogAction(ctx, "player_stats_auth_fail", apiKey, 0, endpoint, nil, false, err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Get player stats
	stats, err := h.mcp.GetPlayerStats(ctx, userID, playerID)
	if err != nil {
		h.auditLogger.LogAction(ctx, "player_stats_error", apiKey, userID, endpoint, map[string]interface{}{"player_id": playerID}, false, err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Log success
	h.auditLogger.LogAction(ctx, "player_stats_success", apiKey, userID, endpoint, map[string]interface{}{"player_id": playerID}, true, "")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// SubmitAction handles agent-submitted actions
func (h *MCPRESTHandler) SubmitAction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	apiKey := h.getAPIKey(r)
	endpoint := "/mcp/action"

	// Apply rate limiting
	allowed, err := h.filter.Allow(apiKey)
	if !allowed {
		h.auditLogger.LogAction(ctx, "action_denied", apiKey, 0, endpoint, nil, false, err.Error())
		http.Error(w, err.Error(), http.StatusTooManyRequests)
		return
	}

	// Authenticate
	userID, err := h.authenticate(ctx, apiKey)
	if err != nil {
		h.auditLogger.LogAction(ctx, "action_auth_fail", apiKey, 0, endpoint, nil, false, err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Parse request body
	var req mcp.SubmitActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.auditLogger.LogAction(ctx, "action_parse_error", apiKey, userID, endpoint, nil, false, err.Error())
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Log request
	h.auditLogger.LogAction(ctx, "action_request", apiKey, userID, endpoint, map[string]interface{}{
		"action_type": req.ActionType,
		"match_id":    req.MatchID,
		"fighter_id":  req.FighterID,
	}, true, "")

	// Submit action
	result, err := h.mcp.SubmitAction(ctx, userID, req)
	if err != nil {
		h.auditLogger.LogAction(ctx, "action_error", apiKey, userID, endpoint, map[string]interface{}{
			"action_type": req.ActionType,
		}, false, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Log success/failure
	successStr := "failure"
	if result.Success {
		successStr = "success"
	}
	h.auditLogger.LogAction(ctx, "action_"+successStr, apiKey, userID, endpoint, map[string]interface{}{
		"action_type": req.ActionType,
	}, result.Success, result.Message)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
