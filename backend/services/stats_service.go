package services

import (
	"strings"
)

type CollectionStats struct {
	MonsterCount int            `json:"monster"`
	SpellCount   int            `json:"spell"`
	TrapCount    int            `json:"trap"`
	Attributes   map[string]int `json:"attributes"`
}

// CalculateCollectionStats computes card type distribution for a user's collection.
func CalculateCollectionStats(userID uint) (CollectionStats, error) {
	userCards, err := GetCollectionByUserID(userID)
	if err != nil {
		return CollectionStats{}, err
	}

	stats := CollectionStats{
		Attributes: make(map[string]int),
	}

	for _, uc := range userCards {
		cardType := uc.Card.Type
		quantity := uc.Quantity

		switch {
		case containsIgnoreCase(cardType, "monster"):
			stats.MonsterCount += quantity

			// Check MonsterCard details
			if uc.Card.MonsterCard != nil {
				attr := strings.ToUpper(uc.Card.MonsterCard.Attribute)
				if attr != "" {
					stats.Attributes[attr] += quantity
				}
			}
		case containsIgnoreCase(cardType, "spell"):
			stats.SpellCount += quantity
		case containsIgnoreCase(cardType, "trap"):
			stats.TrapCount += quantity
		}
	}

	return stats, nil
}

// Helper function to match type regardless of case
func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
