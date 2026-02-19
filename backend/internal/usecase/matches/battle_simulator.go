package matches

import (
	"encoding/json"
	"math"
	"math/rand"
	"sort"
	"time"

	"empoweredpixels/internal/domain/combat"
	"empoweredpixels/internal/domain/roster"
)

type BattleSimulator struct {
	rng *rand.Rand
}

func NewBattleSimulator() *BattleSimulator {
	return &BattleSimulator{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

type BattleOptions struct {
	MaxRounds int
	MapSize   float64
}

const (
	MaxActionGauge = 100.0
)

func (s *BattleSimulator) Run(matchID string, fighters []roster.Fighter, options BattleOptions) (*combat.MatchResult, error) {
	if options.MaxRounds == 0 {
		options.MaxRounds = 1000 // Increased for tick-based epic battles
	}
	if options.MapSize == 0 {
		options.MapSize = 30.0
	}

	entities := s.initializeEntities(fighters, options.MapSize)
	scores := make(map[string]*combat.FighterScore)
	for _, e := range entities {
		scores[e.ID] = &combat.FighterScore{FighterID: e.ID}
	}

	var roundTicks []combat.RoundTick

	// Initial Spawn
	spawnTicks := s.generateSpawnTicks(entities)
	roundTicks = append(roundTicks, combat.RoundTick{Round: 0, Ticks: spawnTicks})

	// Battle Loop (Tick-based)
	for tickCounter := 1; tickCounter <= options.MaxRounds; tickCounter++ {
		var currentTickEvents []combat.Tick
		alive := s.getAlive(entities)

		// Termination check: check if more than one team is still alive
		teamsAlive := make(map[string]bool)
		for _, e := range alive {
			team := "NO_TEAM"
			if e.TeamID != nil {
				team = *e.TeamID
			}
			teamsAlive[team] = true
		}

		if len(teamsAlive) <= 1 {
			break
		}

		// Increase Action Gauge based on Agility
		for _, e := range alive {
			// Base growth (buffed to 20) + Agility scaling
			growth := 20.0 + (float64(e.Stats.Agility) / 1.5)
			e.ActionGauge += growth
		}

		// Check if battle should end due to no alive entities
		if len(alive) == 0 {
			break
		}

		// Process actors who reached the threshold
		// Sort by overflow to handle simultaneous turns fairly
		sort.Slice(alive, func(i, j int) bool {
			return alive[i].ActionGauge > alive[j].ActionGauge
		})

		for _, attacker := range alive {
			if attacker.ActionGauge < MaxActionGauge {
				continue
			}

			// Attacker takes turn
			if attacker.CurrentHP <= 0 {
				attacker.ActionGauge = 0 // Reset even if dead to prevent loops
				continue
			}

			target := s.findNearestTarget(attacker, alive)
			if target == nil {
				attacker.ActionGauge = 0
				continue
			}

			dist := s.distance(attacker, target)
			skill := s.selectSkill(attacker)

			if dist <= skill.Range() {
				eventTicks, err := skill.Execute(attacker, target)
				if err == nil {
					currentTickEvents = append(currentTickEvents, eventTicks...)
					for _, t := range eventTicks {
						if t.Type == "died" {
							scores[attacker.ID].Kills++
							scores[target.ID].Deaths++
						}
					}
				}
			} else {
				currentTickEvents = append(currentTickEvents, s.moveTowards(attacker, target)...)
			}

			// Turn cost
			attacker.ActionGauge -= MaxActionGauge
			if attacker.ActionGauge < 0 {
				attacker.ActionGauge = 0
			}
		}

		// Check if all entities died during this tick - end immediately
		alive = s.getAlive(entities)
		if len(alive) == 0 {
			if len(currentTickEvents) > 0 {
				roundTicks = append(roundTicks, combat.RoundTick{Round: tickCounter, Ticks: currentTickEvents})
			}
			break
		}

		if len(currentTickEvents) > 0 {
			roundTicks = append(roundTicks, combat.RoundTick{Round: tickCounter, Ticks: currentTickEvents})
		}
	}

	return &combat.MatchResult{
		MatchID:    matchID,
		RoundTicks: roundTicks,
		Scores:     s.finalizeScores(scores),
	}, nil
}

func (s *BattleSimulator) initializeEntities(fighters []roster.Fighter, mapSize float64) []*combat.Entity {
	entities := make([]*combat.Entity, len(fighters))
	for i, f := range fighters {
		maxHP := 200 + (f.Vitality * 20) // Significantly buffed HP for longer battles
		entities[i] = &combat.Entity{
			ID:           f.ID,
			Name:         f.Name,
			Level:        f.Level,
			MaxHP:        maxHP,
			CurrentHP:    maxHP,
			AttunementID: f.AttunementID,
			TeamID:       f.TeamID,
			X:            s.rng.Float64() * mapSize,
			Y:            s.rng.Float64() * mapSize,
			Stats: combat.Stats{
				Power:          f.Power,
				ConditionPower: f.ConditionPower,
				Precision:      f.Precision,
				Ferocity:       f.Ferocity,
				Accuracy:       f.Accuracy,
				Agility:        f.Agility,
				Armor:          f.Armor,
				Vitality:       f.Vitality,
				ParryChance:    f.ParryChance,
				HealingPower:   f.HealingPower,
				Speed:          f.Speed,
				Vision:         f.Vision,
			},
			Momentum:    1.0, // Start with neutral momentum
			ActionGauge: s.rng.Float64() * 20.0,
		}
	}
	return entities
}

func (s *BattleSimulator) generateSpawnTicks(entities []*combat.Entity) []combat.Tick {
	ticks := make([]combat.Tick, len(entities))
	for i, e := range entities {
		spawn := combat.EventSpawn{FighterID: e.ID, X: e.X, Y: e.Y, HP: e.CurrentHP}
		p, _ := json.Marshal(spawn)
		ticks[i] = combat.Tick{Type: "spawn", Payload: p}
	}
	return ticks
}

func (s *BattleSimulator) getAlive(entities []*combat.Entity) []*combat.Entity {
	var alive []*combat.Entity
	for _, e := range entities {
		if e.CurrentHP > 0 {
			alive = append(alive, e)
		}
	}
	return alive
}

func (s *BattleSimulator) sortByInitiative(entities []*combat.Entity) {
	sort.Slice(entities, func(i, j int) bool {
		initI := entities[i].Stats.Speed + entities[i].Stats.Agility + s.rng.Intn(10)
		initJ := entities[j].Stats.Speed + entities[j].Stats.Agility + s.rng.Intn(10)
		return initI > initJ
	})
}

func (s *BattleSimulator) findNearestTarget(attacker *combat.Entity, alive []*combat.Entity) *combat.Entity {
	var nearest *combat.Entity
	minDist := math.MaxFloat64
	for _, e := range alive {
		if e.ID == attacker.ID {
			continue
		}

		// Don't attack teammates
		if attacker.TeamID != nil && e.TeamID != nil && *attacker.TeamID == *e.TeamID {
			continue
		}

		d := s.distance(attacker, e)
		if d < minDist {
			minDist = d
			nearest = e
		}
	}
	return nearest
}

func (s *BattleSimulator) distance(a, b *combat.Entity) float64 {
	return math.Sqrt(math.Pow(a.X-b.X, 2) + math.Pow(a.Y-b.Y, 2))
}

func (s *BattleSimulator) selectSkill(e *combat.Entity) combat.Skill {
	// Simple logic: pick weapon based on highest stat
	weaponID := "Greatsword"
	if e.Stats.Precision > e.Stats.Power {
		weaponID = "Bow"
	} else if e.Stats.Agility > e.Stats.Power {
		weaponID = "Dagger"
	}

	skills := combat.GetSkillsByWeapon(weaponID)
	if len(skills) > 0 {
		return skills[0]
	}
	return combat.NewGreatswordBlow()
}

func (s *BattleSimulator) moveTowards(attacker *combat.Entity, target *combat.Entity) []combat.Tick {
	fromX, fromY := attacker.X, attacker.Y

	// Speed-based movement distance
	moveDist := 3.0 + (float64(attacker.Stats.Speed) / 8.0)

	angle := math.Atan2(target.Y-attacker.Y, target.X-attacker.X)
	attacker.X += math.Cos(angle) * moveDist
	attacker.Y += math.Sin(angle) * moveDist

	// Reset combo/loss momentum on move (repositioning)
	attacker.Combo = 0
	attacker.Momentum -= 0.1
	if attacker.Momentum < 0 {
		attacker.Momentum = 0
	}

	move := combat.EventMove{
		FighterID: attacker.ID,
		FromX:     fromX, FromY: fromY,
		ToX: attacker.X, ToY: attacker.Y,
	}
	p, _ := json.Marshal(move)
	return []combat.Tick{{Type: "move", Payload: p}}
}

func (s *BattleSimulator) finalizeScores(scoreMap map[string]*combat.FighterScore) []combat.FighterScore {
	scores := make([]combat.FighterScore, 0, len(scoreMap))
	for _, s := range scoreMap {
		scores = append(scores, *s)
	}
	return scores
}
