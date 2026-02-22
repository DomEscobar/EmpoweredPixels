package leagues

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	domainleagues "empoweredpixels/internal/domain/leagues"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService is a mock implementation of the leagues Service interface for testing
type MockService struct {
	mock.Mock
}

func (m *MockService) List(ctx context.Context) ([]domainleagues.League, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domainleagues.League), args.Error(1)
}

func (m *MockService) Get(ctx context.Context, id int) (*domainleagues.League, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainleagues.League), args.Error(1)
}

func (m *MockService) Subscribe(ctx context.Context, userID int64, leagueID int, fighterID string) error {
	args := m.Called(ctx, userID, leagueID, fighterID)
	return args.Error(0)
}

func (m *MockService) Unsubscribe(ctx context.Context, userID int64, leagueID int, fighterID string) error {
	args := m.Called(ctx, userID, leagueID, fighterID)
	return args.Error(0)
}

func (m *MockService) Subscriptions(ctx context.Context, leagueID int) ([]domainleagues.LeagueSubscription, error) {
	args := m.Called(ctx, leagueID)
	return args.Get(0).([]domainleagues.LeagueSubscription), args.Error(1)
}

func (m *MockService) SubscriptionsForUser(ctx context.Context, leagueID int, userID int64) ([]domainleagues.LeagueSubscription, error) {
	args := m.Called(ctx, leagueID, userID)
	return args.Get(0).([]domainleagues.LeagueSubscription), args.Error(1)
}

func (m *MockService) Matches(ctx context.Context, leagueID int, page int, pageSize int) ([]domainleagues.LeagueMatch, int, error) {
	args := m.Called(ctx, leagueID, page, pageSize)
	return args.Get(0).([]domainleagues.LeagueMatch), args.Get(1).(int), args.Error(2)
}

func (m *MockService) GetHighScores(ctx context.Context, leagueID int, lastMatches int) ([]domainleagues.LeagueHighscore, error) {
	args := m.Called(ctx, leagueID, lastMatches)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domainleagues.LeagueHighscore), args.Error(1)
}

func (m *MockService) GetLastWinner(ctx context.Context, leagueID int) (*domainleagues.LeagueWinner, error) {
	args := m.Called(ctx, leagueID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainleagues.LeagueWinner), args.Error(1)
}

func (m *MockService) CreateLeague(ctx context.Context, name string, options []byte, isDeactivated bool) (*domainleagues.League, error) {
	args := m.Called(ctx, name, options, isDeactivated)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainleagues.League), args.Error(1)
}

func (m *MockService) UpdateLeague(ctx context.Context, id int, name string, options []byte, isDeactivated bool) (*domainleagues.League, error) {
	args := m.Called(ctx, id, name, options, isDeactivated)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainleagues.League), args.Error(1)
}

func (m *MockService) DeleteLeague(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestHandler_Highscores_DefaultLastMatches(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	leagueID := 1
	mockService.On("GetHighScores", mock.Anything, leagueID, 50).Return([]domainleagues.LeagueHighscore{
		{FighterID: "f1", FighterName: "Fighter1", Username: "user1", Score: 100},
	}, nil)

	req := httptest.NewRequest("GET", "/api/league/1/highscores", nil)
	w := httptest.NewRecorder()

	handler.Highscores(w, req, strconv.Itoa(leagueID))

	assert.Equal(t, http.StatusOK, w.Code)
	var result []leagueHighscoreDto
	err := json.NewDecoder(w.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "f1", result[0].FighterID)
	assert.Equal(t, "Fighter1", result[0].FighterName)
	assert.Equal(t, "user1", result[0].Username)
	assert.Equal(t, 100, result[0].Score)

	mockService.AssertExpectations(t)
}

func TestHandler_Highscores_CustomLastMatches(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	leagueID := 2
	lastMatches := 10
	mockService.On("GetHighScores", mock.Anything, leagueID, lastMatches).Return([]domainleagues.LeagueHighscore{
		{FighterID: "f2", FighterName: "Fighter2", Username: "user2", Score: 50},
	}, nil)

	req := httptest.NewRequest("GET", "/api/league/2/highscores?lastMatches=10", nil)
	w := httptest.NewRecorder()

	handler.Highscores(w, req, strconv.Itoa(leagueID))

	assert.Equal(t, http.StatusOK, w.Code)
	var result []leagueHighscoreDto
	err := json.NewDecoder(w.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	mockService.AssertExpectations(t)
}

func TestHandler_Highscores_InvalidLastMatches(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	leagueID := 3
	// Invalid lastMatches should default to 50
	mockService.On("GetHighScores", mock.Anything, leagueID, 50).Return([]domainleagues.LeagueHighscore{}, nil)

	req := httptest.NewRequest("GET", "/api/league/3/highscores?lastMatches=abc", nil)
	w := httptest.NewRecorder()

	handler.Highscores(w, req, strconv.Itoa(leagueID))

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandler_Highscores_EmptyResult(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	leagueID := 4
	mockService.On("GetHighScores", mock.Anything, leagueID, 50).Return([]domainleagues.LeagueHighscore{}, nil)

	req := httptest.NewRequest("GET", "/api/league/4/highscores", nil)
	w := httptest.NewRecorder()

	handler.Highscores(w, req, strconv.Itoa(leagueID))

	assert.Equal(t, http.StatusOK, w.Code)
	var result []leagueHighscoreDto
	err := json.NewDecoder(w.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Len(t, result, 0)
	mockService.AssertExpectations(t)
}

func TestHandler_Highscores_ServiceError(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	leagueID := 5
	mockService.On("GetHighScores", mock.Anything, leagueID, 50).Return([]domainleagues.LeagueHighscore{}, assert.AnError)

	req := httptest.NewRequest("GET", "/api/league/5/highscores", nil)
	w := httptest.NewRecorder()

	handler.Highscores(w, req, strconv.Itoa(leagueID))

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandler_Highscores_InvalidLeagueID(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	req := httptest.NewRequest("GET", "/api/league/invalid/highscores", nil)
	w := httptest.NewRecorder()

	handler.Highscores(w, req, "invalid")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	// service should not be called
	mockService.AssertNotCalled(t, "GetHighScores")
}

func TestHandler_Matches_ReturnsCorrectTotalCount(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	leagueID := 1
	page := 2
	pageSize := 5
	totalMatches := 42 // total in DB

	matches := []domainleagues.LeagueMatch{
		{LeagueID: leagueID, MatchID: "m1", Started: nil},
		{LeagueID: leagueID, MatchID: "m2", Started: nil},
		{LeagueID: leagueID, MatchID: "m3", Started: nil},
		{LeagueID: leagueID, MatchID: "m4", Started: nil},
		{LeagueID: leagueID, MatchID: "m5", Started: nil},
	}

	mockService.On("Matches", mock.Anything, leagueID, page, pageSize).Return(matches, totalMatches, nil)

	reqBody := `{"page": 2, "pageSize": 5}`
	req := httptest.NewRequest("POST", "/api/league/1/matches", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Matches(w, req, strconv.Itoa(leagueID))

	assert.Equal(t, http.StatusOK, w.Code)
	var resp pageDto[leagueMatchDto]
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, page, resp.Page)
	assert.Equal(t, pageSize, resp.PageSize)
	assert.Equal(t, totalMatches, resp.TotalCount) // totalCount should reflect DB total, not len(items)
	assert.Len(t, resp.Items, len(matches))

	mockService.AssertExpectations(t)
}

func TestHandler_Matches_TotalCountUnaffectedByPagination(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	leagueID := 1
	totalMatches := 100

	// Request page 1 with size 10
	matchesPage1 := make([]domainleagues.LeagueMatch, 10)
	for i := 0; i < 10; i++ {
		matchesPage1[i] = domainleagues.LeagueMatch{LeagueID: leagueID, MatchID: fmt.Sprintf("m%d", i+1), Started: nil}
	}
	mockService.On("Matches", mock.Anything, leagueID, 1, 10).Return(matchesPage1, totalMatches, nil)

	// Request page 2 with size 10
	matchesPage2 := make([]domainleagues.LeagueMatch, 10)
	for i := 0; i < 10; i++ {
		matchesPage2[i] = domainleagues.LeagueMatch{LeagueID: leagueID, MatchID: fmt.Sprintf("m%d", i+11), Started: nil}
	}
	mockService.On("Matches", mock.Anything, leagueID, 2, 10).Return(matchesPage2, totalMatches, nil)

	// Page 1
	reqBody1 := `{"page": 1, "pageSize": 10}`
	req1 := httptest.NewRequest("POST", "/api/league/1/matches", strings.NewReader(reqBody1))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	handler.Matches(w1, req1, strconv.Itoa(leagueID))
	var resp1 pageDto[leagueMatchDto]
	json.NewDecoder(w1.Body).Decode(&resp1)

	assert.Equal(t, totalMatches, resp1.TotalCount)
	assert.Len(t, resp1.Items, 10)

	// Page 2
	reqBody2 := `{"page": 2, "pageSize": 10}`
	req2 := httptest.NewRequest("POST", "/api/league/1/matches", strings.NewReader(reqBody2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	handler.Matches(w2, req2, strconv.Itoa(leagueID))
	var resp2 pageDto[leagueMatchDto]
	json.NewDecoder(w2.Body).Decode(&resp2)

	assert.Equal(t, totalMatches, resp2.TotalCount)
	assert.Len(t, resp2.Items, 10)

	mockService.AssertExpectations(t)
}

func TestHandler_Matches_ServiceError(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	leagueID := 1
	mockService.On("Matches", mock.Anything, leagueID, 1, 20).Return([]domainleagues.LeagueMatch{}, 0, assert.AnError)

	reqBody := `{"page": 1, "pageSize": 20}`
	req := httptest.NewRequest("POST", "/api/league/1/matches", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Matches(w, req, strconv.Itoa(leagueID))

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandler_Matches_DefaultPaging(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	leagueID := 1
	expectedPage := 1
	expectedPageSize := 20
	matches := []domainleagues.LeagueMatch{{LeagueID: leagueID, MatchID: "m1", Started: nil}}
	totalMatches := 5

	mockService.On("Matches", mock.Anything, leagueID, expectedPage, expectedPageSize).Return(matches, totalMatches, nil)

	// Empty body should default to page=1, pageSize=20
	req := httptest.NewRequest("POST", "/api/league/1/matches", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Matches(w, req, strconv.Itoa(leagueID))

	assert.Equal(t, http.StatusOK, w.Code)
	var resp pageDto[leagueMatchDto]
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, expectedPage, resp.Page)
	assert.Equal(t, expectedPageSize, resp.PageSize)
	assert.Equal(t, totalMatches, resp.TotalCount)
	mockService.AssertExpectations(t)
}

func TestHandler_LastWinner_Success(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	leagueID := 1
	winner := &domainleagues.LeagueWinner{
		FighterID:   "f123",
		FighterName: "Thunder",
		Username:    "player1",
		MatchID:     "match-abc",
	}

	mockService.On("GetLastWinner", mock.Anything, leagueID).Return(winner, nil)

	req := httptest.NewRequest("GET", "/api/league/1/winner", nil)
	w := httptest.NewRecorder()

	handler.LastWinner(w, req, strconv.Itoa(leagueID))

	assert.Equal(t, http.StatusOK, w.Code)
	var resp domainleagues.LeagueWinner
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, winner.FighterID, resp.FighterID)
	assert.Equal(t, winner.FighterName, resp.FighterName)
	assert.Equal(t, winner.Username, resp.Username)
	assert.Equal(t, winner.MatchID, resp.MatchID)

	mockService.AssertExpectations(t)
}

func TestHandler_LastWinner_NoContent(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	leagueID := 2
	mockService.On("GetLastWinner", mock.Anything, leagueID).Return(nil, nil)

	req := httptest.NewRequest("GET", "/api/league/2/winner", nil)
	w := httptest.NewRecorder()

	handler.LastWinner(w, req, strconv.Itoa(leagueID))

	assert.Equal(t, http.StatusNoContent, w.Code)
	// Body should be empty
	assert.Equal(t, "", w.Body.String())
	mockService.AssertExpectations(t)
}

func TestHandler_LastWinner_ServiceError(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	leagueID := 3
	mockService.On("GetLastWinner", mock.Anything, leagueID).Return(nil, assert.AnError)

	req := httptest.NewRequest("GET", "/api/league/3/winner", nil)
	w := httptest.NewRecorder()

	handler.LastWinner(w, req, strconv.Itoa(leagueID))

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandler_LastWinner_InvalidLeagueID(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	req := httptest.NewRequest("GET", "/api/league/invalid/winner", nil)
	w := httptest.NewRecorder()

	handler.LastWinner(w, req, "invalid")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	// service should not be called
	mockService.AssertNotCalled(t, "GetLastWinner")
}
