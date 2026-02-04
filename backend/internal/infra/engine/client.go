package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type Client struct {
	baseURL string
	http    *http.Client
}

type Config struct {
	BaseURL string
	Timeout time.Duration
}

type MatchRequest struct {
	MatchID  string          `json:"matchId"`
	Options  json.RawMessage `json:"options"`
	Fighters []FighterInput  `json:"fighters"`
}

type FighterInput struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type MatchResult struct {
	RoundTicks json.RawMessage `json:"roundTicks"`
	Scores     []ScoreResult   `json:"scores"`
}

type ScoreResult struct {
	FighterID    string `json:"fighterId"`
	TotalKills   int    `json:"totalKills"`
	TotalDeaths  int    `json:"totalDeaths"`
	TotalAssists int    `json:"totalAssists"`
}

func NewClient(cfg Config) *Client {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	return &Client{
		baseURL: cfg.BaseURL,
		http:    &http.Client{Timeout: timeout},
	}
}

func (c *Client) RunMatch(ctx context.Context, req MatchRequest) (*MatchResult, error) {
	if c.baseURL == "" {
		return nil, errors.New("engine url not configured")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/match/run", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New("engine error")
	}

	var result MatchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
