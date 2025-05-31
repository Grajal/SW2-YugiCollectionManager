package services

import (
	"strings"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
)

type AvgStats struct {
	AvgATK float64 `json:"avg_atk"`
	AvgDEF float64 `json:"avg_def"`
}

type CollectionStats struct {
	MonsterCount int            `json:"monster"`
	SpellCount   int            `json:"spell"`
	TrapCount    int            `json:"trap"`
	Attributes   map[string]int `json:"attributes"`
	AverageStats AvgStats       `json:"average_stats"`
}

// CardWithQuantity is a generic wrapper for any card + quantity structure,
// allowing reuse of stat functions across collections and decks.
type CardWithQuantity struct {
	Card     models.Card
	Quantity int
}

// CalculateDeckStats computes statistics from a list of DeckCards,
// including type distribution, monster attributes, and average ATK/DEF.
func CalculateDeckStats(deckCards []models.DeckCard) CollectionStats {
	cards := convertDeckCardsToGeneric(deckCards)

	monsterCount, spellCount, trapCount := countCardTypes(cards)
	attributes := countMonsterAttributes(cards)
	avgStats := computeAverageStats(cards)

	return CollectionStats{
		MonsterCount: monsterCount,
		SpellCount:   spellCount,
		TrapCount:    trapCount,
		Attributes:   attributes,
		AverageStats: avgStats,
	}
}

// convertDeckCardsToGeneric converts a slice of DeckCard into []CardWithQuantity
func convertDeckCardsToGeneric(deckCards []models.DeckCard) []CardWithQuantity {
	result := make([]CardWithQuantity, 0, len(deckCards))
	for _, dc := range deckCards {
		result = append(result, CardWithQuantity{
			Card:     dc.Card,
			Quantity: dc.Quantity,
		})
	}
	return result
}

// CalculateCollectionStats returns overall statistics for a user's card collection,
// including counts of monsters, spells, and traps, monster attribute distribution,
// and average ATK/DEF for monster cards.
func CalculateCollectionStats(userID uint) (CollectionStats, error) {
	userCards, err := GetCollectionByUserID(userID)
	if err != nil {
		return CollectionStats{}, err
	}

	cards := convertUserCardsToGeneric(userCards)

	monsterCount, spellCount, trapCount := countCardTypes(cards)
	attributes := countMonsterAttributes(cards)
	avgStats := computeAverageStats(cards)

	return CollectionStats{
		MonsterCount: monsterCount,
		SpellCount:   spellCount,
		TrapCount:    trapCount,
		Attributes:   attributes,
		AverageStats: avgStats,
	}, nil
}

// convertUserCardsToGeneric transforms []UserCard into []CardWithQuantity
func convertUserCardsToGeneric(userCards []models.UserCard) []CardWithQuantity {
	result := make([]CardWithQuantity, 0, len(userCards))
	for _, uc := range userCards {
		result = append(result, CardWithQuantity{
			Card:     uc.Card,
			Quantity: uc.Quantity,
		})
	}
	return result
}

// countCardTypes calculates the number of Monster, Spell, and Trap cards from a generic card slice.
func countCardTypes(cards []CardWithQuantity) (int, int, int) {
	var monsterCount, spellCount, trapCount int

	for _, c := range cards {
		qty := c.Quantity
		cardType := c.Card.Type

		switch {
		case containsIgnoreCase(cardType, "monster"):
			monsterCount += qty
		case containsIgnoreCase(cardType, "spell"):
			spellCount += qty
		case containsIgnoreCase(cardType, "trap"):
			trapCount += qty
		}
	}

	return monsterCount, spellCount, trapCount
}

// countMonsterAttributes tallies the distribution of monster attributes (e.g., DARK, LIGHT).
func countMonsterAttributes(cards []CardWithQuantity) map[string]int {
	attributes := make(map[string]int)

	for _, c := range cards {
		if c.Card.MonsterCard != nil {
			attr := strings.ToUpper(c.Card.MonsterCard.Attribute)
			if attr != "" {
				attributes[attr] += c.Quantity
			}
		}
	}

	return attributes
}

// computeAverageStats calculates average ATK/DEF across all Monster cards, weighted by quantity.
func computeAverageStats(cards []CardWithQuantity) AvgStats {
	var totalATK, totalDEF, count int

	for _, c := range cards {
		qty := c.Quantity
		if c.Card.MonsterCard != nil {
			atk := c.Card.MonsterCard.Atk
			def := c.Card.MonsterCard.Def

			if atk >= 0 {
				totalATK += atk * qty
			}
			if def >= 0 {
				totalDEF += def * qty
			}
			count += qty
		}
	}

	if count == 0 {
		return AvgStats{}
	}

	return AvgStats{
		AvgATK: float64(totalATK) / float64(count),
		AvgDEF: float64(totalDEF) / float64(count),
	}
}

// containsIgnoreCase checks if s contains substr (case-insensitive).
func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
