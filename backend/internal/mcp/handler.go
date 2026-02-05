package mcp

import (
	"context"
	"errors"
	"fmt"

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
