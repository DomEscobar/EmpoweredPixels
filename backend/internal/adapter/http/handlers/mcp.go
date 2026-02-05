package handlers

import (
	"encoding/json"
	"net/http"

	"empoweredpixels/internal/mcp"
)

type MCPHandler struct {
	mcp *mcp.MCPHandler
}

func NewMCPHandler(mcp *mcp.MCPHandler) *MCPHandler {
	return &MCPHandler{mcp: mcp}
}

type CallToolRequest struct {
	Tool string                 `json:"tool"`
	Args map[string]interface{} `json:"args"`
}

func (h *MCPHandler) Call(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("X-EP-AI-KEY")
	if apiKey == "" {
		http.Error(w, "missing X-EP-AI-KEY", http.StatusUnauthorized)
		return
	}

	var req CallToolRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.mcp.CallTool(r.Context(), apiKey, req.Tool, req.Args)
	if err != nil {
		http.Error(w, err.Error(), http.StatusTooManyRequests)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
