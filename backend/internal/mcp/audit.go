package mcp

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

// AuditLogEntry represents a single audit log entry
type AuditLogEntry struct {
	Timestamp   time.Time              `json:"timestamp"`
	Action      string                 `json:"action"`
	AgentID     string                 `json:"agent_id"`
	UserID      int64                  `json:"user_id"`
	Endpoint    string                 `json:"endpoint"`
	RequestBody map[string]interface{} `json:"request_body,omitempty"`
	Success     bool                   `json:"success"`
	ErrorMsg    string                 `json:"error_msg,omitempty"`
	ClientIP    string                 `json:"client_ip,omitempty"`
}

// AuditLogger handles MCP audit logging
type AuditLogger struct {
	mu       sync.Mutex
	logger   *log.Logger
	file     *os.File
	enabled  bool
	entries  []AuditLogEntry // In-memory buffer for recent entries
	maxSize  int
}

// NewAuditLogger creates a new audit logger
func NewAuditLogger(logPath string) (*AuditLogger, error) {
	al := &AuditLogger{
		enabled: true,
		maxSize: 1000,
		entries: make([]AuditLogEntry, 0, 1000),
	}

	if logPath != "" {
		file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		al.file = file
		al.logger = log.New(file, "", 0)
	} else {
		al.logger = log.New(os.Stdout, "[MCP-AUDIT] ", log.LstdFlags)
	}

	return al, nil
}

// Log records an audit log entry
func (al *AuditLogger) Log(ctx context.Context, entry AuditLogEntry) {
	if al == nil || !al.enabled {
		return
	}

	al.mu.Lock()
	defer al.mu.Unlock()

	entry.Timestamp = time.Now().UTC()

	// Write to log
	data, err := json.Marshal(entry)
	if err == nil {
		al.logger.Println(string(data))
	}

	// Buffer in memory (keep last 1000)
	al.entries = append(al.entries, entry)
	if len(al.entries) > al.maxSize {
		al.entries = al.entries[len(al.entries)-al.maxSize:]
	}
}

// GetRecent returns recent audit log entries
func (al *AuditLogger) GetRecent(count int) []AuditLogEntry {
	al.mu.Lock()
	defer al.mu.Unlock()

	if count > len(al.entries) {
		count = len(al.entries)
	}
	
	result := make([]AuditLogEntry, count)
	start := len(al.entries) - count
	copy(result, al.entries[start:])
	return result
}

// Close closes the audit log file
func (al *AuditLogger) Close() error {
	if al.file != nil {
		return al.file.Close()
	}
	return nil
}

// LogAction is a helper to log an action with common fields
func (al *AuditLogger) LogAction(ctx context.Context, action, agentID string, userID int64, endpoint string, body map[string]interface{}, success bool, errMsg string) {
	entry := AuditLogEntry{
		Action:      action,
		AgentID:     agentID,
		UserID:      userID,
		Endpoint:    endpoint,
		RequestBody: body,
		Success:     success,
		ErrorMsg:    errMsg,
	}
	al.Log(ctx, entry)
}
