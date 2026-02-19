package roster

import (
	"context"
	"errors"
	"testing"

	"empoweredpixels/internal/domain/roster"
	"github.com/google/uuid"
)

type MockSquadRepository struct {
	squadName        string
	squads           map[string]*roster.Squad
	activeSquad      map[int64]*roster.Squad
	shouldDeactivate error
	shouldCreate     error
	shouldGetActive  error
}

func NewMockSquadRepository() *MockSquadRepository {
	return &MockSquadRepository{
		squads:           make(map[string]*roster.Squad),
		activeSquad:      make(map[int64]*roster.Squad),
		shouldDeactivate: nil,
		shouldCreate:     nil,
		shouldGetActive:  nil,
	}
}

func (m *MockSquadRepository) Create(ctx context.Context, squad *roster.Squad) error {
	if m.shouldCreate != nil {
		return m.shouldCreate
	}
	m.squads[squad.ID] = squad
	m.activeSquad[squad.UserID] = squad
	return nil
}

func (m *MockSquadRepository) GetActiveByUserID(ctx context.Context, userID int64) (*roster.Squad, error) {
	if m.shouldGetActive != nil {
		return nil, m.shouldGetActive
	}
	return m.activeSquad[userID], nil
}

func (m *MockSquadRepository) DeactivateAll(ctx context.Context, userID int64) error {
	if m.shouldDeactivate != nil {
		return m.shouldDeactivate
	}
	delete(m.activeSquad, userID)
	return nil
}

type MockFighterRepository struct {
	fighters map[string]*roster.Fighter
}

func NewMockFighterRepository() *MockFighterRepository {
	return &MockFighterRepository{
		fighters: make(map[string]*roster.Fighter),
	}
}

func (m *MockFighterRepository) ListByUser(ctx context.Context, userID int64) ([]roster.Fighter, error) {
	return nil, nil
}
func (m *MockFighterRepository) GetByUserAndID(ctx context.Context, userID int64, id string) (*roster.Fighter, error) {
	return nil, nil
}
func (m *MockFighterRepository) GetByID(ctx context.Context, id string) (*roster.Fighter, error) {
	return m.fighters[id], nil
}
func (m *MockFighterRepository) NameExists(ctx context.Context, name string) (bool, error) {
	return false, nil
}
func (m *MockFighterRepository) UserHasFighter(ctx context.Context, userID int64) (bool, error) {
	return false, nil
}
func (m *MockFighterRepository) CountByUser(ctx context.Context, userID int64) (int, error) {
	return 0, nil
}
func (m *MockFighterRepository) Create(ctx context.Context, fighter *roster.Fighter) error {
	return nil
}
func (m *MockFighterRepository) Update(ctx context.Context, fighter *roster.Fighter) error {
	return nil
}
func (m *MockFighterRepository) SoftDelete(ctx context.Context, userID int64, id string) error {
	return nil
}

func TestSetActiveSquad(t *testing.T) {
	ctx := context.Background()
	fighterRepo := NewMockFighterRepository()

	tests := []struct {
		name       string
		userID     int64
		squadName  string
		fighterIDs []string
		repo       *MockSquadRepository
		wantErr    bool
		errorMsg   string
		verify     func(*testing.T, *roster.Squad)
	}{
		{
			name:      "Create active squad with 3 fighters",
			userID:    1,
			squadName: "Test Squad",
			fighterIDs: []string{
				uuid.NewString(),
				uuid.NewString(),
				uuid.NewString(),
			},
			repo: NewMockSquadRepository(),
			verify: func(t *testing.T, squad *roster.Squad) {
				if squad == nil {
					t.Fatal("Expected squad to not be nil")
				}
				if squad.Name != "Test Squad" {
					t.Errorf("Expected squad name 'Test Squad', got '%s'", squad.Name)
				}
				if squad.UserID != 1 {
					t.Errorf("Expected userID 1, got %d", squad.UserID)
				}
				if squad.IsActive != true {
					t.Error("Expected squad to be active")
				}
				if len(squad.Members) != 3 {
					t.Errorf("Expected 3 members, got %d", len(squad.Members))
				}
			},
		},
		{
			name:      "Limit to 3 fighters",
			userID:    2,
			squadName: "Large Squad",
			fighterIDs: []string{
				uuid.NewString(),
				uuid.NewString(),
				uuid.NewString(),
				uuid.NewString(),
				uuid.NewString(),
			},
			repo: NewMockSquadRepository(),
			verify: func(t *testing.T, squad *roster.Squad) {
				if len(squad.Members) != 3 {
					t.Errorf("Expected max 3 fighters, got %d", len(squad.Members))
				}
			},
		},
		{
			name:       "Create squad with 1 fighter",
			userID:     3,
			squadName:  "Small Squad",
			fighterIDs: []string{uuid.NewString()},
			repo:       NewMockSquadRepository(),
			verify: func(t *testing.T, squad *roster.Squad) {
				if len(squad.Members) != 1 {
					t.Errorf("Expected 1 member, got %d", len(squad.Members))
				}
			},
		},
		{
			name:       "Handle deactivate all error",
			userID:     4,
			squadName:  "Error Squad",
			fighterIDs: []string{uuid.NewString(), uuid.NewString()},
			repo: func() *MockSquadRepository {
				m := NewMockSquadRepository()
				m.shouldDeactivate = errors.New("database error")
				return m
			}(),
			wantErr:  true,
			errorMsg: "database error",
		},
		{
			name:       "Handle create error",
			userID:     5,
			squadName:  "Create Error Squad",
			fighterIDs: []string{uuid.NewString()},
			repo: func() *MockSquadRepository {
				m := NewMockSquadRepository()
				m.shouldCreate = errors.New("duplicate squad")
				return m
			}(),
			wantErr:  true,
			errorMsg: "duplicate squad",
		},
		{
			name:       "Handle get active error",
			userID:     6,
			squadName:  "Get Error Squad",
			fighterIDs: []string{uuid.NewString()},
			repo: func() *MockSquadRepository {
				m := NewMockSquadRepository()
				m.shouldGetActive = errors.New("query failed")
				return m
			}(),
			wantErr:  true,
			errorMsg: "query failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewSquadService(tt.repo, fighterRepo)

			squad, err := service.SetActiveSquad(ctx, tt.userID, tt.squadName, tt.fighterIDs)

			if (err != nil) != tt.wantErr {
				t.Errorf("SetActiveSquad() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Fatal("Expected error but got nil")
				}
				if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
				return
			}

			if squad == nil {
				t.Fatal("Expected squad to not be nil")
			}

			if tt.verify != nil {
				tt.verify(t, squad)
			}
		})
	}
}

func TestGetActiveSquad(t *testing.T) {
	ctx := context.Background()
	fighterRepo := NewMockFighterRepository()

	tests := []struct {
		name      string
		userID    int64
		repo      *MockSquadRepository
		wantSquad bool
		wantErr   bool
	}{
		{
			name:   "Get existing active squad",
			userID: 1,
			repo: func() *MockSquadRepository {
				m := NewMockSquadRepository()
				m.activeSquad[1] = &roster.Squad{
					ID:       uuid.NewString(),
					UserID:   1,
					Name:     "Active Squad",
					IsActive: true,
					Members:  []roster.Member{{FighterID: "fighter1", SlotIndex: 0}},
				}
				return m
			}(),
			wantSquad: true,
		},
		{
			name:      "Get non-existent squad",
			userID:    999,
			repo:      NewMockSquadRepository(),
			wantSquad: false,
		},
		{
			name:   "Handle repository error",
			userID: 2,
			repo: func() *MockSquadRepository {
				m := NewMockSquadRepository()
				m.shouldGetActive = errors.New("repository error")
				return m
			}(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewSquadService(tt.repo, fighterRepo)

			squad, err := service.GetActiveSquad(ctx, tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetActiveSquad() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Fatal("Expected error but got nil")
				}
				return
			}

			if tt.wantSquad && squad == nil {
				t.Fatal("Expected squad to not be nil")
			}

			if !tt.wantSquad && squad != nil {
				t.Error("Expected squad to be nil")
			}
		})
	}
}

func TestGetActiveSquadParallelHydration(t *testing.T) {
	ctx := context.Background()
	fighterRepo := NewMockFighterRepository()

	// Setup fighters in mock repo
	fighterRepo.fighters["f1"] = &roster.Fighter{ID: "f1", Name: "Fighter 1"}
	fighterRepo.fighters["f2"] = &roster.Fighter{ID: "f2", Name: "Fighter 2"}
	fighterRepo.fighters["f3"] = &roster.Fighter{ID: "f3", Name: "Fighter 3"}

	repo := NewMockSquadRepository()
	repo.activeSquad[1] = &roster.Squad{
		ID:       uuid.NewString(),
		UserID:   1,
		Name:     "Test Squad",
		IsActive: true,
		Members: []roster.Member{
			{FighterID: "f1", SlotIndex: 0},
			{FighterID: "f2", SlotIndex: 1},
			{FighterID: "f3", SlotIndex: 2},
		},
	}

	service := NewSquadService(repo, fighterRepo)
	squad, err := service.GetActiveSquad(ctx, 1)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if squad == nil {
		t.Fatal("Expected squad to not be nil")
	}

	if len(squad.Members) != 3 {
		t.Fatalf("Expected 3 members, got %d", len(squad.Members))
	}

	for i, member := range squad.Members {
		if member.Fighter == nil {
			t.Errorf("Member at index %d has nil Fighter data", i)
		} else if member.Fighter.ID != member.FighterID {
			t.Errorf("Member at index %d has mismatching Fighter ID: got %s, want %s", i, member.Fighter.ID, member.FighterID)
		}
	}
}

func TestNewSquadService(t *testing.T) {
	repo := NewMockSquadRepository()
	fighterRepo := NewMockFighterRepository()
	service := NewSquadService(repo, fighterRepo)

	if service == nil {
		t.Fatal("Expected service to not be nil")
	}

	if service.repo != repo {
		t.Error("Expected repo to be set correctly")
	}
}
