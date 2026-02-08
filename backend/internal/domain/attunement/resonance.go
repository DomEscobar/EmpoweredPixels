package attunement

import (
	"fmt"
	"sort"
	"strings"
)

// ResonancePattern represents the attunement pattern of a squad
type ResonancePattern struct {
	PrimaryElements     []Element
	HarmonicPairs       []ElementPair
	DissonantPairs      []ElementPair
	HarmonyScore        int
	TierName            string
	AuraColor           string
}

// ElementPair represents two elements that interact
type ElementPair struct {
	Element1 Element
	Element2 Element
	Type     string // "harmonic" or "dissonant"
}

// HarmonicPairs define elements that strengthen each other
var HarmonicPairs = []ElementPair{
	{Fire, Air, "harmonic"},
	{Air, Fire, "harmonic"},
	{Water, Earth, "harmonic"},
	{Earth, Water, "harmonic"},
	{Light, Dark, "harmonic"},
	{Dark, Light, "harmonic"},
}

// DissonantPairs define elements that weaken each other
var DissonantPairs = []ElementPair{
	{Fire, Water, "dissonant"},
	{Water, Fire, "dissonant"},
	{Air, Earth, "dissonant"},
	{Earth, Air, "dissonant"},
}

// CalculateHarmony calculates the squad harmony score based on fighter attunements.
// Harmonic pairs: Fire↔Air, Water↔Earth, Light↔Dark (+8-12 points each)
// Dissonant pairs: Fire↔Water, Air↔Earth (-5 penalty)
// Score formula: (harmonic_count × 10) + (dissonant_count × -5) + base 25
func CalculateHarmony(elements []Element) (score int, pattern *ResonancePattern) {
	if len(elements) == 0 {
		return 0, &ResonancePattern{
			PrimaryElements: []Element{},
			HarmonicPairs:   []ElementPair{},
			DissonantPairs:  []ElementPair{},
			HarmonyScore:    0,
			TierName:        "Discordant",
			AuraColor:       "#808080",
		}
	}

	// Start with base score
	score = 25

	harmonicCount := 0
	dissonantCount := 0
	harmonicPairs := []ElementPair{}
	dissonantPairs := []ElementPair{}

	// Check all pairwise combinations
	for i := 0; i < len(elements); i++ {
		for j := i + 1; j < len(elements); j++ {
			pair := ElementPair{Element1: elements[i], Element2: elements[j]}

			// Check if harmonic
			if isHarmonicPair(elements[i], elements[j]) {
				harmonicCount++
				harmonicPairs = append(harmonicPairs, pair)
				score += 10
			}

			// Check if dissonant
			if isDissonantPair(elements[i], elements[j]) {
				dissonantCount++
				dissonantPairs = append(dissonantPairs, pair)
				score -= 5
			}
		}
	}

	// Clamp score to 0-100
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	// Determine tier name and aura color
	tierName := getTierName(score)
	auraColor := getAuraColor(score)

	pattern = &ResonancePattern{
		PrimaryElements: elements,
		HarmonicPairs:   harmonicPairs,
		DissonantPairs:  dissonantPairs,
		HarmonyScore:    score,
		TierName:        tierName,
		AuraColor:       auraColor,
	}

	return score, pattern
}

// isHarmonicPair checks if two elements form a harmonic pair
func isHarmonicPair(e1, e2 Element) bool {
	switch {
	case (e1 == Fire && e2 == Air) || (e1 == Air && e2 == Fire):
		return true
	case (e1 == Water && e2 == Earth) || (e1 == Earth && e2 == Water):
		return true
	case (e1 == Light && e2 == Dark) || (e1 == Dark && e2 == Light):
		return true
	}
	return false
}

// isDissonantPair checks if two elements form a dissonant pair
func isDissonantPair(e1, e2 Element) bool {
	switch {
	case (e1 == Fire && e2 == Water) || (e1 == Water && e2 == Fire):
		return true
	case (e1 == Air && e2 == Earth) || (e1 == Earth && e2 == Air):
		return true
	}
	return false
}

// getTierName returns the tier name based on harmony score
func getTierName(score int) string {
	switch {
	case score >= 76:
		return "Resonant"
	case score >= 51:
		return "Harmonized"
	case score >= 26:
		return "Aligned"
	default:
		return "Discordant"
	}
}

// getAuraColor returns the aura color based on harmony score
func getAuraColor(score int) string {
	switch {
	case score >= 76:
		// Resonant: plasma/rainbow gradient (simulated with vibrant color)
		return "#FFD700" // Gold, but would be animated in UI
	case score >= 51:
		// Harmonized: gold
		return "#FFD700"
	case score >= 26:
		// Aligned: blue
		return "#4169E1"
	default:
		// Discordant: gray
		return "#808080"
	}
}

// GetBonusMultiplier returns damage and defense multipliers based on harmony score
func GetBonusMultiplier(score int) (dmgMultiplier, defMultiplier float64) {
	switch {
	case score >= 76:
		return 1.12, 1.06
	case score >= 51:
		return 1.08, 1.04
	case score >= 26:
		return 1.0, 1.0
	default:
		return 0.95, 0.95 // Debuff for dissonant squads
	}
}

// GetDissonanceWarning returns a warning message if the pattern has dissonant pairs
func GetDissonanceWarning(pattern *ResonancePattern) string {
	if len(pattern.DissonantPairs) == 0 {
		return ""
	}

	var warnings []string
	for _, pair := range pattern.DissonantPairs {
		warnings = append(warnings, fmt.Sprintf("%s and %s are in conflict", pair.Element1, pair.Element2))
	}

	return "⚠️ Dissonance Detected: " + strings.Join(warnings, "; ")
}

// FormatElementList returns a formatted list of elements
func FormatElementList(elements []Element) []string {
	var formatted []string
	for _, e := range elements {
		formatted = append(formatted, strings.ToTitle(strings.ToLower(string(e))))
	}
	sort.Strings(formatted)
	return formatted
}
