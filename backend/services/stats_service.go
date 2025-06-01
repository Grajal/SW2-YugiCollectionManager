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

// CalculateCollectionStats returns overall statistics for a user's card collection,
// including counts of monsters, spells, and traps, monster attribute distribution,
// and average ATK/DEF for monster cards.
func CalculateCollectionStats(userID uint, collectionService CollectionService) (CollectionStats, error) {
	userCards, err := collectionService.GetUserCollection(userID)
	if err != nil {
		return CollectionStats{}, err
	}

	monsterCount, spellCount, trapCount := countCardTypes(userCards)
	attributes := countMonsterAttributes(userCards)
	avgStats := computeAverageStats(userCards)

	return CollectionStats{
		MonsterCount: monsterCount,
		SpellCount:   spellCount,
		TrapCount:    trapCount,
		Attributes:   attributes,
		AverageStats: avgStats,
	}, nil
}

// countCardTypes calculates the total number of Monster, Spell, and Trap cards
// in the user's collection based on card type and quantity.
func countCardTypes(userCards []models.UserCard) (int, int, int) {
	var monsterCount, spellCount, trapCount int

	for _, uc := range userCards {
		qty := uc.Quantity
		cardType := uc.Card.Type

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

// countMonsterAttributes returns a map of monster attributes (e.g., DARK, LIGHT)
// and their respective total quantities in the user's collection.
func countMonsterAttributes(userCards []models.UserCard) map[string]int {
	attributes := make(map[string]int)

	for _, uc := range userCards {
		if uc.Card.MonsterCard != nil {
			attr := strings.ToUpper(uc.Card.MonsterCard.Attribute)
			if attr != "" {
				attributes[attr] += uc.Quantity
			}
		}
	}

	return attributes
}

// computeAverageStats calculates the average ATK and DEF values for all Monster cards
// in the user's collection, weighted by the quantity of each card.
func computeAverageStats(userCards []models.UserCard) AvgStats {
	var totalATK, totalDEF, count int

	for _, uc := range userCards {
		qty := uc.Quantity
		if uc.Card.MonsterCard != nil {
			atk := uc.Card.MonsterCard.Atk
			def := uc.Card.MonsterCard.Def

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

// Helper function to match type regardless of case
func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
