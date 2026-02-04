package combat

import (
	"encoding/json"
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
	
	// Apply armor reduction (simple formula for now)
	damage = damage - (target.Stats.Armor / 10)
	if damage < 1 {
		damage = 1
	}

	target.CurrentHP -= damage
	if target.CurrentHP < 0 {
		target.CurrentHP = 0
	}

	attackEvent := EventAttack{
		AttackerID: attacker.ID,
		TargetID:   target.ID,
		SkillID:    s.id,
		Damage:     damage,
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
