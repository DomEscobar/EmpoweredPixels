package combat

import (
	"encoding/json"
	"math/rand"
)

type Attunement interface {
	ID() string
	StrongAgainst() string
	WeakAgainst() string
	OnAttack(attacker *Entity, target *Entity) []Tick
}

type BaseAttunement struct {
	id            string
	strongAgainst string
	weakAgainst   string
}

func (a *BaseAttunement) ID() string            { return a.id }
func (a *BaseAttunement) StrongAgainst() string { return a.strongAgainst }
func (a *BaseAttunement) WeakAgainst() string   { return a.weakAgainst }

type FireAttunement struct{ BaseAttunement }

func (a *FireAttunement) OnAttack(attacker *Entity, target *Entity) []Tick {
	if rand.Float64() > 0.1 { // 10% chance
		return nil
	}
	// Apply Burn (dummy condition for now)
	event := map[string]string{
		"fighterId": attacker.ID,
		"targetId":  target.ID,
		"condition": "Burn",
	}
	p, _ := json.Marshal(event)
	return []Tick{{Type: "conditionApplied", Payload: p}}
}

func NewFireAttunement() Attunement {
	return &FireAttunement{BaseAttunement{
		id:            "Fire",
		strongAgainst: "Wind",
		weakAgainst:   "Water",
	}}
}

// ... other attunements can be added similarly
