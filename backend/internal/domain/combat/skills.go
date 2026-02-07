package combat

import (
	"encoding/json"
	"math"
	"math/rand"
)

type Skill interface {
	ID() string
	Name() string
	Range() float64
	Execute(attacker *Entity, target *Entity) ([]Tick, error)
}

type BaseDamageSkill struct {
	id        string
	name      string
	minDamage int
	maxDamage int
	rng       float64
}

func (s *BaseDamageSkill) ID() string      { return s.id }
func (s *BaseDamageSkill) Name() string    { return s.name }
func (s *BaseDamageSkill) Range() float64 { return s.rng }

func (s *BaseDamageSkill) Execute(attacker *Entity, target *Entity) ([]Tick, error) {
	damage := rand.Intn(s.maxDamage-s.minDamage+1) + s.minDamage

	// Apply Combo & Momentum Multipliers
	// Each combo point adds 5% damage
	// Each 1.0 momentum adds 10% damage
	comboBonus := 1.0 + (float64(attacker.Combo) * 0.05)
	momentumBonus := 1.0 + (attacker.Momentum * 0.1)
	
	finalDamage := float64(damage) * comboBonus * momentumBonus
	damage = int(math.Floor(finalDamage))
	
	// Increment Combo and Momentum on hit
	attacker.Combo++
	attacker.Momentum += 0.2
	if attacker.Momentum > 5.0 {
		attacker.Momentum = 5.0
	}

	// Apply armor reduction (simple formula for now)
	damage = damage - (target.Stats.Armor / 10)
	if damage < 1 {
		damage = 1
	}

	target.CurrentHP -= damage
	// Target loses momentum when hit
	target.Momentum -= 0.5
	if target.Momentum < 0 {
		target.Momentum = 0
	}
	// Target combo breaks when hit
	target.Combo = 0

	if target.CurrentHP < 0 {
		target.CurrentHP = 0
	}

	attackEvent := EventAttack{
		AttackerID: attacker.ID,
		TargetID:   target.ID,
		SkillID:    s.id,
		Damage:     damage,
		Combo:      attacker.Combo,
		Momentum:   attacker.Momentum,
	}
	
	payload, _ := json.Marshal(attackEvent)
	ticks := []Tick{{Type: "attack", Payload: payload}}

	if target.CurrentHP <= 0 {
		diedEvent := EventDied{
			FighterID: target.ID,
			KillerID:  attacker.ID,
		}
		diedPayload, _ := json.Marshal(diedEvent)
		ticks = append(ticks, Tick{Type: "died", Payload: diedPayload})
	}

	return ticks, nil
}

func NewBowShot() Skill {
	return &BaseDamageSkill{
		id:        "27CCB150-9EA2-4FD1-8DFC-EB13E512C225",
		name:      "Shot",
		minDamage: 4,
		maxDamage: 6,
		rng:       20.0,
	}
}

func NewDaggerSlice() Skill {
	return &BaseDamageSkill{
		id:        "DaggerSlice",
		name:      "Slice",
		minDamage: 3,
		maxDamage: 5,
		rng:       2.0,
	}
}

func NewGlaiveSwing() Skill {
	return &BaseDamageSkill{
		id:        "GlaiveSwing",
		name:      "Swing",
		minDamage: 5,
		maxDamage: 8,
		rng:       4.0,
	}
}

func NewGreatswordBlow() Skill {
	return &BaseDamageSkill{
		id:        "GreatswordBlow",
		name:      "Blow",
		minDamage: 8,
		maxDamage: 12,
		rng:       3.0,
	}
}

func GetSkillsByWeapon(weaponID string) []Skill {
	switch weaponID {
	case "Bow":
		return []Skill{NewBowShot()}
	case "Dagger":
		return []Skill{NewDaggerSlice()}
	case "Glaive":
		return []Skill{NewGlaiveSwing()}
	case "Greatsword":
		return []Skill{NewGreatswordBlow()}
	default:
		return []Skill{NewGreatswordBlow()}
	}
}
