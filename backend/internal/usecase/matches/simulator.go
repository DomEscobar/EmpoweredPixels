package matches

import (
	"encoding/json"
	"math"
	"math/rand"

	"empoweredpixels/internal/domain/combat"
	"empoweredpixels/internal/domain/inventory"
	"empoweredpixels/internal/domain/roster"
	"github.com/google/uuid"
)

type Simulator struct {
}

func NewSimulator() *Simulator {
	return &Simulator{}
}

func (s *Simulator) Run(matchID string, fighters []roster.Fighter, fighterEquipment map[string][]inventory.Equipment, options MatchOptions) (*combat.MatchResult, error) {
	entities := make([]*combat.Entity, 0, len(fighters)+(func() int {
		if options.BotCount != nil {
			return *options.BotCount
		}
		return 0
	}()))
	scores := make(map[string]*combat.FighterScore)

	for _, f := range fighters {
		maxHP := 100 + (f.Vitality * 10)
		entities = append(entities, &combat.Entity{
			ID:           f.ID,
			Name:         f.Name,
			Level:        f.Level,
			MaxHP:        maxHP,
			CurrentHP:    maxHP,
			AttunementID: f.AttunementID,
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
		})
		scores[f.ID] = &combat.FighterScore{FighterID: f.ID}
	}

	// Add bots
	if options.BotCount != nil {
		level := 10
		if options.BotPowerlevel != nil {
			level = *options.BotPowerlevel
		}
		for i := 0; i < *options.BotCount; i++ {
			id := uuid.NewString()
			name := "Bot " + id[:4]
			maxHP := 100 + (level * 5)
			entities = append(entities, &combat.Entity{
				ID:        id,
				Name:      name,
				Level:     level,
				MaxHP:     maxHP,
				CurrentHP: maxHP,
				Stats: combat.Stats{
					Power:    level,
					Accuracy: level,
					Speed:    level,
					Armor:    level / 2,
				},
			})
			scores[id] = &combat.FighterScore{FighterID: id}
		}
	}

	// Initial spawn events
	var roundTicks []combat.RoundTick
	var spawnTicks []combat.Tick
	for _, e := range entities {
		e.X = rand.Float64() * 20
		e.Y = rand.Float64() * 20
		spawn := combat.EventSpawn{FighterID: e.ID, X: e.X, Y: e.Y, HP: e.CurrentHP}
		p, _ := json.Marshal(spawn)
		spawnTicks = append(spawnTicks, combat.Tick{Type: "spawn", Payload: p})
	}
	roundTicks = append(roundTicks, combat.RoundTick{Round: 0, Ticks: spawnTicks})

	// Run rounds
	for round := 1; round <= 50; round++ {
		var ticks []combat.Tick
		alive := getAlive(entities)
		if len(alive) <= 1 {
			break
		}

		// Shuffle turn order
		rand.Shuffle(len(alive), func(i, j int) { alive[i], alive[j] = alive[j], alive[i] })

		for _, attacker := range alive {
			if attacker.CurrentHP <= 0 {
				continue
			}

			target := findNearestTarget(attacker, alive)
			if target == nil {
				continue
			}

			dist := distance(attacker, target)
			skill := getSkillForAttacker(attacker)
			attunement := getAttunement(attacker.AttunementID)

			if dist <= skill.Range() {
				// Attack
				eventTicks, _ := skill.Execute(attacker, target)
				ticks = append(ticks, eventTicks...)

				// Apply Attunement effect
				if attunement != nil {
					attunementTicks := attunement.OnAttack(attacker, target)
					ticks = append(ticks, attunementTicks...)
				}

				// Update scores
				for _, t := range eventTicks {
					if t.Type == "died" {
						scores[attacker.ID].Kills++
						scores[target.ID].Deaths++
					}
				}
			} else {
				// Move towards target
				attacker.Momentum -= 0.1
				if attacker.Momentum < 0 {
					attacker.Momentum = 0
				}
				attacker.Combo = 0

				fromX, fromY := attacker.X, attacker.Y
				moveDist := 2.0 + (float64(attacker.Stats.Speed) / 10.0)
				angle := math.Atan2(target.Y-attacker.Y, target.X-attacker.X)
				attacker.X += math.Cos(angle) * moveDist
				attacker.Y += math.Sin(angle) * moveDist

				move := combat.EventMove{
					FighterID: attacker.ID,
					FromX:     fromX, FromY: fromY,
					ToX: attacker.X, ToY: attacker.Y,
				}
				p, _ := json.Marshal(move)
				ticks = append(ticks, combat.Tick{Type: "move", Payload: p})
			}
		}

		if len(ticks) > 0 {
			roundTicks = append(roundTicks, combat.RoundTick{Round: round, Ticks: ticks})
		}
	}

	finalScores := make([]combat.FighterScore, 0, len(scores))
	for _, s := range scores {
		finalScores = append(finalScores, *s)
	}

	return &combat.MatchResult{
		MatchID:    matchID,
		RoundTicks: roundTicks,
		Scores:     finalScores,
	}, nil
}

func getAlive(entities []*combat.Entity) []*combat.Entity {
	var alive []*combat.Entity
	for _, e := range entities {
		if e.CurrentHP > 0 {
			alive = append(alive, e)
		}
	}
	return alive
}

func findNearestTarget(attacker *combat.Entity, alive []*combat.Entity) *combat.Entity {
	var nearest *combat.Entity
	minDist := math.MaxFloat64
	for _, e := range alive {
		if e.ID == attacker.ID {
			continue
		}
		d := distance(attacker, e)
		if d < minDist {
			minDist = d
			nearest = e
		}
	}
	return nearest
}

func distance(a, b *combat.Entity) float64 {
	return math.Sqrt(math.Pow(a.X-b.X, 2) + math.Pow(a.Y-b.Y, 2))
}

func getSkillForAttacker(e *combat.Entity) combat.Skill {
	weaponID := "Greatsword"
	if e.Stats.Speed > 15 {
		weaponID = "Bow"
	} else if e.Stats.Agility > 15 {
		weaponID = "Dagger"
	}

	skills := combat.GetSkillsByWeapon(weaponID)
	if len(skills) > 0 {
		return skills[0]
	}
	return combat.NewGreatswordBlow()
}

func getAttunement(id *string) combat.Attunement {
	if id == nil {
		return nil
	}
	switch *id {
	case "Fire":
		return combat.NewFireAttunement()
	default:
		return nil
	}
}
