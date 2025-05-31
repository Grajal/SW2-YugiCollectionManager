package services

import (
	"strings"
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

func CalculateCollectionStats(userID uint) (CollectionStats, error) {
	userCards, err := GetCollectionByUserID(userID)
	if err != nil {
		return CollectionStats{}, err
	}

	stats := CollectionStats{
		Attributes: make(map[string]int),
	}

	var (
		totalATK    int
		totalDEF    int
		atkDefCount int
	)

	for _, uc := range userCards {
		cardType := uc.Card.Type
		quantity := uc.Quantity

		switch {
		case containsIgnoreCase(cardType, "monster"):
			stats.MonsterCount += quantity

			// Monster-specific logic
			if uc.Card.MonsterCard != nil {
				attr := strings.ToUpper(uc.Card.MonsterCard.Attribute)
				if attr != "" {
					stats.Attributes[attr] += quantity
				}

				atk := uc.Card.MonsterCard.Atk
				def := uc.Card.MonsterCard.Def

				if atk >= 0 {
					totalATK += atk * quantity
				}
				if def >= 0 {
					totalDEF += def * quantity
				}

				atkDefCount += quantity
			}
		case containsIgnoreCase(cardType, "spell"):
			stats.SpellCount += quantity
		case containsIgnoreCase(cardType, "trap"):
			stats.TrapCount += quantity
		}
	}

	if atkDefCount > 0 {
		stats.AverageStats.AvgATK = float64(totalATK) / float64(atkDefCount)
		stats.AverageStats.AvgDEF = float64(totalDEF) / float64(atkDefCount)
	}

	return stats, nil
}

// Helper function to match type regardless of case
func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
