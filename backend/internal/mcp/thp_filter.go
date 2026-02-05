package mcp

import (
	"errors"
	"sync"
	"time"
)

// FairnessFilter manages the "The Human Pace" (THP) limits
type FairnessFilter struct {
	mu            sync.Mutex
	lastActions   map[string][]time.Time
	limitRequests int
	windowSize    time.Duration
}

// NewFairnessFilter creates a new rate limiter with specified limits
func NewFairnessFilter(limit int, window time.Duration) *FairnessFilter {
	return &FairnessFilter{
		lastActions:   make(map[string][]time.Time),
		limitRequests: limit,
		windowSize:    window,
	}
}

// DefaultFairnessFilter creates a rate limiter with 100 requests per minute
func DefaultFairnessFilter() *FairnessFilter {
	return NewFairnessFilter(100, 1*time.Minute)
}

// Allow checks if the AI-Key is within the APM limits
func (f *FairnessFilter) Allow(aiKey string) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	now := time.Now()
	actions, ok := f.lastActions[aiKey]
	if !ok {
		f.lastActions[aiKey] = []time.Time{now}
		return true, nil
	}

	// Clean up old timestamps
	var validActions []time.Time
	for _, t := range actions {
		if now.Sub(t) < f.windowSize {
			validActions = append(validActions, t)
		}
	}

	if len(validActions) >= f.limitRequests {
		f.lastActions[aiKey] = validActions
		return false, errors.New("rate limit exceeded: The Human Pace (THP) protection active")
	}

	validActions = append(validActions, now)
	f.lastActions[aiKey] = validActions
	return true, nil
}
