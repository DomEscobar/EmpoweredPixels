package mcp

import (
	"context"
	"errors"
	"fmt"
	"time"

	"empoweredpixels/internal/domain/inventory"
	"empoweredpixels/internal/domain/matches"
	"empoweredpixels/internal/domain/roster"
	"empoweredpixels/internal/usecase/identity"
	inventoryusecase "empoweredpixels/internal/usecase/inventory"
	leaguesusecase "empoweredpixels/internal/usecase/leagues"
	matchesusecase "empoweredpixels/internal/usecase/matches"
	"empoweredpixels/internal/usecase/rewards"
	rosterusecase "empoweredpixels/internal/usecase/roster"
)

// MCPHandler handles the core logic for the Model Context Protocol server
type MCPHandler struct {
	filter           *FairnessFilter
	identityService  *identity.Service
	rosterService    *rosterusecase.Service
	inventoryService inventoryusecase.Service
	leagueService    *leaguesusecase.Service
	matchService     *matchesusecase.Service
	rewardService    *rewards.Service
}

func NewMCPHandler(
	filter *FairnessFilter,
	identityService *identity.Service,
	rosterService *rosterusecase.Service,
	inventoryService inventoryusecase.Service,
	leagueService *leaguesusecase.Service,
	matchService *matchesusecase.Service,
	rewardService *rewards.Service,
) *MCPHandler {
	return &MCPHandler{
		filter:           filter,
		identityService:  identityService,
		rosterService:    rosterService,
		inventoryService: inventoryService,
		leagueService:    leagueService,
		matchService:     matchService,
		rewardService:    rewardService,
	}
}

// Authenticate validates the X-EP-AI-KEY and returns the associated user context
func (h *MCPHandler) Authenticate(ctx context.Context, apiKey string) (int64, error) {
	if apiKey == "" {
		return 0, errors.New("missing X-EP-AI-KEY")
	}

	// For Alpha Test, we map a specific key to our AI-Alpha user (ID 4)
	if apiKey == "AI-ALPHA-TEST-KEY-2026" {
		return 4, nil
	}

	return 0, errors.New("invalid AI Key")
}

// CallTool executes a tool while applying the Fairness Filter
func (h *MCPHandler) CallTool(ctx context.Context, apiKey string, toolName string, args map[string]interface{}) (interface{}, error) {
	// First: Apply Fairness Filter (THP)
	allowed, err := h.filter.Allow(apiKey)
	if !allowed {
		return nil, err
	}

	userID, err := h.Authenticate(ctx, apiKey)
	if err != nil {
		return nil, err
	}

	switch toolName {
	case "get_roster":
		return h.handleGetRoster(ctx, userID)
	case "create_fighter":
		return h.handleCreateFighter(ctx, userID, args)
	case "get_inventory":
		return h.handleGetInventory(ctx, userID)
	case "equip_fighter":
		return h.handleEquipFighter(ctx, userID, args)
	case "join_league":
		return h.handleJoinLeague(ctx, userID, args)
	case "browse_matches":
		return h.handleBrowseMatches(ctx)
	case "create_match":
		return h.handleCreateMatch(ctx, userID, args)
	case "join_match":
		return h.handleJoinMatch(ctx, userID, args)
	case "issue_starter_loot":
		return h.handleIssueStarterLoot(ctx, userID)
	default:
		return nil, fmt.Errorf("unknown tool: %s", toolName)
	}
}

func (h *MCPHandler) handleGetRoster(ctx context.Context, userID int64) (interface{}, error) {
	return h.rosterService.List(ctx, userID)
}

func (h *MCPHandler) handleCreateFighter(ctx context.Context, userID int64, args map[string]interface{}) (interface{}, error) {
	name, ok := args["name"].(string)
	if !ok {
		return nil, errors.New("missing name")
	}
	return h.rosterService.Create(ctx, userID, name)
}

func (h *MCPHandler) handleGetInventory(ctx context.Context, userID int64) (interface{}, error) {
	return h.inventoryService.InventoryPage(ctx, userID, 1, 100)
}

func (h *MCPHandler) handleEquipFighter(ctx context.Context, userID int64, args map[string]interface{}) (interface{}, error) {
	fighterID, ok := args["fighter_id"].(string)
	if !ok {
		return nil, errors.New("missing fighter_id")
	}
	equipmentID, ok := args["equipment_id"].(string)
	if !ok {
		return nil, errors.New("missing equipment_id")
	}
	err := h.inventoryService.Equip(ctx, userID, equipmentID, &fighterID)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("Equipment %s equipped to fighter %s", equipmentID, fighterID), nil
}

func (h *MCPHandler) handleJoinLeague(ctx context.Context, userID int64, args map[string]interface{}) (interface{}, error) {
	leagueIDf, ok := args["league_id"].(float64)
	if !ok {
		return nil, errors.New("missing league_id")
	}
	fighterID, ok := args["fighter_id"].(string)
	if !ok {
		return nil, errors.New("missing fighter_id")
	}
	err := h.leagueService.Subscribe(ctx, userID, int(leagueIDf), fighterID)
	if err != nil {
		return nil, err
	}
	return "Joined league successfully", nil
}

func (h *MCPHandler) handleBrowseMatches(ctx context.Context) (interface{}, error) {
	return h.matchService.BrowseMatches(ctx, 1, 20)
}

func (h *MCPHandler) handleCreateMatch(ctx context.Context, userID int64, args map[string]interface{}) (interface{}, error) {
	return h.matchService.CreateMatch(ctx, userID, h.matchService.DefaultOptions())
}

func (h *MCPHandler) handleJoinMatch(ctx context.Context, userID int64, args map[string]interface{}) (interface{}, error) {
	matchID, ok := args["match_id"].(string)
	if !ok {
		return nil, errors.New("missing match_id")
	}
	fighterID, ok := args["fighter_id"].(string)
	if !ok {
		return nil, errors.New("missing fighter_id")
	}
	err := h.matchService.Join(ctx, userID, matchID, fighterID)
	if err != nil {
		return nil, err
	}
	return "Joined match successfully", nil
}

func (h *MCPHandler) handleIssueStarterLoot(ctx context.Context, userID int64) (interface{}, error) {
	// Create some starter equipment for the Alpha Bot via reward system
	_, err := h.rewardService.IssueReward(ctx, userID, "starter_pack")
	if err != nil {
		return nil, err
	}

	// Also claim all immediately
	return h.rewardService.ClaimAll(ctx, userID)
}

// GameState represents the current game state for MCP agents
type GameState struct {
	Matches       []MatchSummary `json:"matches"`
	CurrentMatch  *MatchDetail   `json:"current_match,omitempty"`
	TotalLobbies  int            `json:"total_lobbies"`
	ActiveMatches int            `json:"active_matches"`
}

type MatchSummary struct {
	ID            string    `json:"id"`
	Status        string    `json:"status"`
	Created       time.Time `json:"created"`
	FighterCount  int       `json:"fighter_count"`
	IsJoinable    bool      `json:"is_joinable"`
}

type MatchDetail struct {
	ID            string                `json:"id"`
	Status        string                `json:"status"`
	Created       time.Time             `json:"created"`
	Teams         []matches.MatchTeam   `json:"teams,omitempty"`
	Registrations []matches.MatchRegistration `json:"registrations,omitempty"`
}

// GetGameState returns the current game state (lobbies and active matches)
func (h *MCPHandler) GetGameState(ctx context.Context, userID int64) (*GameState, error) {
	// Get lobby matches
	lobbies, err := h.matchService.BrowseByStatus(ctx, matches.MatchStatusLobby, 1, 20)
	if err != nil {
		return nil, err
	}

	// Get active matches
	active, err := h.matchService.BrowseByStatus(ctx, matches.MatchStatusRunning, 1, 20)
	if err != nil {
		return nil, err
	}

	// Get user's current match
	currentMatch, err := h.matchService.GetCurrentMatch(ctx, userID)
	if err != nil {
		currentMatch = nil
	}

	// Build match summaries
	var summaries []MatchSummary
	for _, m := range lobbies {
		regs, _ := h.matchService.GetRegistrations(ctx, m.ID)
		summaries = append(summaries, MatchSummary{
			ID:           m.ID,
			Status:       m.Status,
			Created:      m.Created,
			FighterCount: len(regs),
			IsJoinable:   m.Status == matches.MatchStatusLobby,
		})
	}

	state := &GameState{
		Matches:       summaries,
		TotalLobbies:  len(lobbies),
		ActiveMatches: len(active),
	}

	// Add current match details if exists
	if currentMatch != nil {
		teams, _ := h.matchService.GetTeams(ctx, currentMatch.ID)
		regs, _ := h.matchService.GetRegistrations(ctx, currentMatch.ID)
		state.CurrentMatch = &MatchDetail{
			ID:            currentMatch.ID,
			Status:        currentMatch.Status,
			Created:       currentMatch.Created,
			Teams:         teams,
			Registrations: regs,
		}
	}

	return state, nil
}

// PlayerStats represents comprehensive player/fighter statistics
type PlayerStats struct {
	FighterID      string                 `json:"fighter_id"`
	Name           string                 `json:"name"`
	Level          int                    `json:"level"`
	Power          int                    `json:"power"`
	Attributes     map[string]int         `json:"attributes"`
	Experience     *roster.FighterExperience `json:"experience,omitempty"`
	Equipment      []inventory.Equipment  `json:"equipment,omitempty"`
	MatchStats     *MatchStatistics       `json:"match_stats,omitempty"`
}

type MatchStatistics struct {
	TotalKills   int `json:"total_kills"`
	TotalDeaths  int `json:"total_deaths"`
	TotalAssists int `json:"total_assists"`
	MatchesPlayed int `json:"matches_played"`
}

// GetPlayerStats returns comprehensive stats for a fighter
func (h *MCPHandler) GetPlayerStats(ctx context.Context, userID int64, fighterID string) (*PlayerStats, error) {
	// Verify fighter belongs to user
	fighter, err := h.rosterService.Get(ctx, userID, fighterID)
	if err != nil {
		return nil, err
	}
	if fighter == nil {
		return nil, errors.New("fighter not found")
	}

	// Get experience
	exp, _ := h.rosterService.GetExperience(ctx, fighterID)

	// Get equipment
	equipment, _ := h.inventoryService.ListByFighter(ctx, userID, fighterID)

	stats := &PlayerStats{
		FighterID:  fighter.ID,
		Name:       fighter.Name,
		Level:      fighter.Level,
		Power:      fighter.Power,
		Attributes: map[string]int{
			"condition_power": fighter.ConditionPower,
			"precision":       fighter.Precision,
			"ferocity":        fighter.Ferocity,
			"accuracy":        fighter.Accuracy,
			"agility":         fighter.Agility,
			"armor":           fighter.Armor,
			"vitality":        fighter.Vitality,
			"parry_chance":    fighter.ParryChance,
			"healing_power":   fighter.HealingPower,
			"speed":           fighter.Speed,
			"vision":          fighter.Vision,
		},
		Experience: exp,
		Equipment:  equipment,
	}

	return stats, nil
}

// SubmitActionRequest represents an action to be submitted by an agent
type SubmitActionRequest struct {
	ActionType string                 `json:"action_type"`
	MatchID    string                 `json:"match_id,omitempty"`
	FighterID  string                 `json:"fighter_id,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// SubmitActionResult represents the result of an action submission
type SubmitActionResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SubmitAction handles agent-submitted actions
func (h *MCPHandler) SubmitAction(ctx context.Context, userID int64, req SubmitActionRequest) (*SubmitActionResult, error) {
	switch req.ActionType {
	case "join_match":
		if req.MatchID == "" || req.FighterID == "" {
			return nil, errors.New("match_id and fighter_id required")
		}
		err := h.matchService.Join(ctx, userID, req.MatchID, req.FighterID)
		if err != nil {
			return &SubmitActionResult{Success: false, Message: err.Error()}, nil
		}
		return &SubmitActionResult{Success: true, Message: "Joined match successfully"}, nil

	case "create_match":
		match, err := h.matchService.CreateMatch(ctx, userID, h.matchService.DefaultOptions())
		if err != nil {
			return &SubmitActionResult{Success: false, Message: err.Error()}, nil
		}
		return &SubmitActionResult{Success: true, Message: "Match created", Data: match}, nil

	case "leave_match":
		if req.MatchID == "" || req.FighterID == "" {
			return nil, errors.New("match_id and fighter_id required")
		}
		err := h.matchService.Leave(ctx, userID, req.MatchID, req.FighterID)
		if err != nil {
			return &SubmitActionResult{Success: false, Message: err.Error()}, nil
		}
		return &SubmitActionResult{Success: true, Message: "Left match"}, nil

	default:
		return nil, fmt.Errorf("unknown action type: %s", req.ActionType)
	}
}
