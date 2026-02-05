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

func (s *Simulator) Run(matchID string, fighters []roster.Fighter, equipment map[string][]inventory.Equipment, options MatchOptions) (*combat.MatchResult, error) {
	entities := make([]*combat.Entity, 0, len(fighters)+(func() int {
		if options.BotCount != nil {
			return *options.BotCount
		}
		return 0
	}()))
	scores := make(map[string]*combat.FighterScore)
	
	// Combo-Momentum state tracking
	comboStates := make(map[string]*combat.ComboMomentumState)

	for _, f := range fighters {
		stats := combat.Stats{
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
		}

		// Apply Equipment Stats
		if items, ok := equipment[f.ID]; ok {
			for _, item := range items {
				applyItemStats(&stats, item)
			}
		}

		maxHP := 100 + (stats.Vitality * 10)
		entities = append(entities, &combat.Entity{
			ID:           f.ID,
			Name:         f.Name,
			Level:        f.Level,
			MaxHP:        maxHP,
			CurrentHP:    maxHP,
			AttunementID: f.AttunementID,
			Stats:        stats,
		})
		scores[f.ID] = &combat.FighterScore{FighterID: f.ID}
		
		// Initialize combo-momentum state
		comboStates[f.ID] = &combat.ComboMomentumState{
			FighterID:       f.ID,
			Momentum:        0,
			ConsecutiveHits: 0,
			CurrentTargetID: "",
			SunderStacks:    0,
			FlurryActive:    false,
			RoundsSinceHit:  0,
		}
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
			
			// Initialize combo-momentum state for bots
			comboStates[id] = &combat.ComboMomentumState{
				FighterID:       id,
				Momentum:        0,
				ConsecutiveHits: 0,
				CurrentTargetID: "",
				SunderStacks:    0,
				FlurryActive:    false,
				RoundsSinceHit:  0,
			}
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

	// Track which entities hit this round for momentum decay
	entitiesThatHitThisRound := make(map[string]bool)

	// Run rounds
	for round := 1; round <= 50; round++ {
		var ticks []combat.Tick
		alive := getAlive(entities)
		if len(alive) <= 1 {
			break
		}

		// Reset hit tracking for this round
		entitiesThatHitThisRound = make(map[string]bool)

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
				// Apply Flurry speed bonus to attack calculation
				attackerComboState := comboStates[attacker.ID]
				if attackerComboState != nil && attackerComboState.FlurryActive {
					// Speed affects accuracy and dodge, apply bonus
					// This is a passive bonus, no visual tick needed
				}

				// Attack
				eventTicks, _ := skill.Execute(attacker, target)
				ticks = append(ticks, eventTicks...)

				// Apply Attunement effect
				if attunement != nil {
					attunementTicks := attunement.OnAttack(attacker, target)
					ticks = append(ticks, attunementTicks...)
				}

				// Check if attack hit (look for attack tick)
				hitTarget := false
				for _, t := range eventTicks {
					if t.Type == "attack" {
						hitTarget = true
						break
					}
				}

				// Combo-Momentum System: Update on successful hit
				if hitTarget && attackerComboState != nil {
					entitiesThatHitThisRound[attacker.ID] = true
					attackerComboState.RoundsSinceHit = 0

					// Update momentum and get events
					momentumTicks, _ := updateMomentumOnHit(attackerComboState, target.ID)
					ticks = append(ticks, momentumTicks...)

					// Apply Sunder debuff to target
					targetComboState := comboStates[target.ID]
					if targetComboState != nil {
						sunderTicks := applySunderDebuff(attackerComboState, targetComboState, target)
						ticks = append(ticks, sunderTicks...)
					}
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
				fromX, fromY := attacker.X, attacker.Y
				moveDist := 2.0 + (float64(attacker.Stats.Speed) / 10.0)
				
				// Apply Flurry speed bonus to movement if active
				attackerComboState := comboStates[attacker.ID]
				if attackerComboState != nil && attackerComboState.FlurryActive {
					moveDist *= 1.0 + (float64(FlurrySpeedBonus) / 100.0)
				}
				
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

		// Apply momentum decay for entities that didn't hit this round
		decayTicks := decayMomentum(comboStates)
		ticks = append(ticks, decayTicks...)

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

func applyItemStats(stats *combat.Stats, item inventory.Equipment) {
	// Base Stats by Item Type (ItemID)
	switch item.ItemID {
	case "sword_01", "iron_sword":
		stats.Power += 5 + item.Enhancement
		stats.Accuracy += 2
	case "greatsword_01":
		stats.Power += 12 + (item.Enhancement * 2)
		stats.Speed -= 2
	case "dagger_01":
		stats.Power += 3 + item.Enhancement
		stats.Agility += 4
		stats.Speed += 2
	case "armor_01", "iron_armor":
		stats.Armor += 8 + item.Enhancement
		stats.Vitality += 2
	case "leather_armor":
		stats.Armor += 3 + item.Enhancement
		stats.Agility += 2
	case "shield_01":
		stats.Armor += 4
		stats.ParryChance += 10 + item.Enhancement
	}

	// Rarity Multiplier
	multiplier := 1.0
	switch item.Rarity {
	case 1: multiplier = 1.1
	case 2: multiplier = 1.25
	case 3: multiplier = 1.5
	case 4: multiplier = 2.0
	}

	if multiplier > 1.0 {
		stats.Power = int(float64(stats.Power) * multiplier)
		stats.Armor = int(float64(stats.Armor) * multiplier)
		stats.Vitality = int(float64(stats.Vitality) * multiplier)
	}
}

// Combo-Momentum System Helpers

const (
	MomentumPerHit      = 10
	MomentumMax         = 100
	MomentumDecay       = 5
	SunderArmorReduction = 0.05 // 5% per stack
	SunderMaxStacks     = 5
	FlurryThreshold     = 50
	FlurrySpeedBonus    = 10 // +10% attack speed
)

// updateMomentumOnHit updates combo state when attacker hits a target
func updateMomentumOnHit(state *combat.ComboMomentumState, targetID string) ([]combat.Tick, bool) {
	var ticks []combat.Tick
	flurryActivated := false

	// Check if hitting same target (combo continues) or different (combo resets)
	if state.CurrentTargetID == targetID {
		// Same target - build combo and momentum
		state.ConsecutiveHits++
		if state.ConsecutiveHits > SunderMaxStacks {
			state.ConsecutiveHits = SunderMaxStacks
		}
		state.Momentum += MomentumPerHit
		if state.Momentum > MomentumMax {
			state.Momentum = MomentumMax
		}
	} else {
		// Different target - reset combo, keep some momentum
		state.ConsecutiveHits = 1
		state.Momentum += MomentumPerHit
		if state.Momentum > MomentumMax {
			state.Momentum = MomentumMax
		}
		state.CurrentTargetID = targetID
		state.SunderStacks = 0 // Reset sunder on target change
	}

	// Check Flurry activation
	wasFlurryActive := state.FlurryActive
	state.FlurryActive = state.Momentum > FlurryThreshold
	if !wasFlurryActive && state.FlurryActive {
		flurryActivated = true
	}

	// Reset rounds since hit
	state.RoundsSinceHit = 0

	// Emit momentum event
	momentumEvent := combat.EventMomentum{
		FighterID:       state.FighterID,
		Momentum:        state.Momentum,
		ConsecutiveHits: state.ConsecutiveHits,
		TargetID:        targetID,
	}
	p, _ := json.Marshal(momentumEvent)
	ticks = append(ticks, combat.Tick{Type: "momentum", Payload: p})

	// Emit flurry event if just activated
	if flurryActivated {
		flurryEvent := combat.EventFlurry{
			FighterID:        state.FighterID,
			AttackSpeedBonus: FlurrySpeedBonus,
		}
		p, _ = json.Marshal(flurryEvent)
		ticks = append(ticks, combat.Tick{Type: "flurry", Payload: p})
	}

	return ticks, flurryActivated
}

// applySunderDebuff applies sunder stacks to target and returns event ticks
func applySunderDebuff(attackerState, targetState *combat.ComboMomentumState, target *combat.Entity) []combat.Tick {
	var ticks []combat.Tick

	if attackerState.ConsecutiveHits <= 1 {
		return ticks // No sunder on first hit
	}

	// Calculate sunder stacks based on consecutive hits
	stacks := attackerState.ConsecutiveHits - 1 // First hit = 0 stacks, 2nd hit = 1 stack, etc.
	if stacks > SunderMaxStacks {
		stacks = SunderMaxStacks
	}

	// Only apply if stacks increased
	if stacks > targetState.SunderStacks {
		oldStacks := targetState.SunderStacks
		targetState.SunderStacks = stacks

		// Calculate armor reduction
		baseArmor := target.Stats.Armor
		totalReduction := float64(baseArmor) * SunderArmorReduction * float64(stacks)
		
		// Calculate newly applied reduction
		newReduction := int(totalReduction - (float64(baseArmor) * SunderArmorReduction * float64(oldStacks)))

		// Apply armor reduction
		target.Stats.Armor = baseArmor - int(totalReduction)

		// Emit sunder event
		sunderEvent := combat.EventSunder{
			TargetID:     target.ID,
			Stacks:       stacks,
			ArmorReduced: newReduction,
		}
		p, _ := json.Marshal(sunderEvent)
		ticks = append(ticks, combat.Tick{Type: "sunder", Payload: p})
	}

	return ticks
}

// decayMomentum applies momentum decay for entities that didn't hit this round
func decayMomentum(states map[string]*combat.ComboMomentumState) []combat.Tick {
	var ticks []combat.Tick

	for _, state := range states {
		state.RoundsSinceHit++
		
		if state.RoundsSinceHit > 0 && state.Momentum > 0 {
			oldMomentum := state.Momentum
			state.Momentum -= MomentumDecay
			if state.Momentum < 0 {
				state.Momentum = 0
			}

			// Check if flurry deactivated
			if oldMomentum > FlurryThreshold && state.Momentum <= FlurryThreshold {
				state.FlurryActive = false
			}

			// Only emit event if momentum changed significantly or hit 0
			if state.Momentum != oldMomentum {
				momentumEvent := combat.EventMomentum{
					FighterID:       state.FighterID,
					Momentum:        state.Momentum,
					ConsecutiveHits: state.ConsecutiveHits,
				}
				p, _ := json.Marshal(momentumEvent)
				ticks = append(ticks, combat.Tick{Type: "momentum_decay", Payload: p})
			}
		}
	}

	return ticks
}

// getEffectiveSpeed returns attack speed with flurry bonus if active
func getEffectiveSpeed(state *combat.ComboMomentumState, baseSpeed int) int {
	if state.FlurryActive {
		bonus := float64(baseSpeed) * (float64(FlurrySpeedBonus) / 100.0)
		return baseSpeed + int(bonus)
	}
	return baseSpeed
}
