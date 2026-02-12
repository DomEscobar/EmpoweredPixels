package leagues

import (
	"context"
	"testing"
	"time"

	"empoweredpixels/internal/domain/leagues"
	"empoweredpixels/internal/domain/roster"
)

type MockLeagueRepo struct{}
func (m *MockLeagueRepo) List(ctx context.Context) ([]leagues.League, error) { return nil, nil }
func (m *MockLeagueRepo) GetByID(ctx context.Context, id int) (*leagues.League, error) {
	return &leagues.League{ID: id, Name: "Test League"}, nil
}

type MockSubscriptionRepo struct {
	createdSub *leagues.LeagueSubscription
}
func (m *MockSubscriptionRepo) ListByLeague(ctx context.Context, leagueID int) ([]leagues.LeagueSubscription, error) { return nil, nil }
func (m *MockSubscriptionRepo) ListByLeagueAndUser(ctx context.Context, leagueID int, userID int64) ([]leagues.LeagueSubscription, error) { return nil, nil }
func (m *MockSubscriptionRepo) Create(ctx context.Context, sub *leagues.LeagueSubscription) error {
	m.createdSub = sub
	return nil
}
func (m *MockSubscriptionRepo) Delete(ctx context.Context, sub *leagues.LeagueSubscription) error { return nil }

type MockFighterRepo struct{}
func (m *MockFighterRepo) GetByUserAndID(ctx context.Context, userID int64, id string) (*roster.Fighter, error) {
	return &roster.Fighter{ID: id, UserID: userID}, nil
}

type MockAchievementRepo struct {
	updatedKey string
	updatedUserID int
}
func (m *MockAchievementRepo) UpdateAchievementProgress(ctx context.Context, userID int, key string, progress int) error {
	m.updatedUserID = userID
	m.updatedKey = key
	return nil
}

func TestSubscribeTriggersAchievement(t *testing.T) {
	leagueRepo := &MockLeagueRepo{}
	subRepo := &MockSubscriptionRepo{}
	fighterRepo := &MockFighterRepo{}
	achieveRepo := &MockAchievementRepo{}

	service := NewService(leagueRepo, subRepo, nil, fighterRepo, achieveRepo, nil)

	err := service.Subscribe(context.Background(), 1, 10, "fighter1")
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	if subRepo.createdSub == nil {
		t.Fatal("Subscription was not created")
	}

	// Wait a bit for the background goroutine
	time.Sleep(50 * time.Millisecond)

	if achieveRepo.updatedKey != "league_enrolled" {
		t.Errorf("Expected achievement key 'league_enrolled', got '%s'", achieveRepo.updatedKey)
	}
	if achieveRepo.updatedUserID != 1 {
		t.Errorf("Expected userID 1, got %d", achieveRepo.updatedUserID)
	}
}
