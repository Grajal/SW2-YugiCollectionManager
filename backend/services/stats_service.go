package services

import (
	"strings"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
)

type AvgStats struct {
	AvgATK float64 `json:"avg_atk"`
	AvgDEF float64 `json:"avg_def"`
}

type Stats struct {
	MonsterCount int            `json:"monster"`
	SpellCount   int            `json:"spell"`
	TrapCount    int            `json:"trap"`
	Attributes   map[string]int `json:"attributes"`
	AverageStats AvgStats       `json:"average_stats"`
	TotalCards   int            `json:"total_cards"`
}

type CardWithQuantity struct {
	Card     models.Card
	Quantity int
}

// CalculateDeckStats computes statistics from a list of DeckCards.
func CalculateDeckStats(deckCards []models.DeckCard) Stats {
	cards := convertDeckCardsToGeneric(deckCards)
	return calculateStats(cards)
}

// CalculateCollectionStats computes statistics for a user's card collection.
func CalculateCollectionStats(userID uint) (Stats, error) {
	userCards, err := GetCollectionByUserID(userID)
	if err != nil {
		return Stats{}, err
	}
	cards := convertUserCardsToGeneric(userCards)
	return calculateStats(cards), nil
}

// Shared logic for computing Stats from generic cards.
func calculateStats(cards []CardWithQuantity) Stats {
	return Stats{
		MonsterCount: countCardTypes(cards, "monster"),
		SpellCount:   countCardTypes(cards, "spell"),
		TrapCount:    countCardTypes(cards, "trap"),
		Attributes:   countMonsterAttributes(cards),
		AverageStats: computeAverageStats(cards),
		TotalCards:   countTotalCards(cards),
	}
}

// Converts []DeckCard into []CardWithQuantity.
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

// Converts []UserCard into []CardWithQuantity.
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

// Counts total number of cards.
func countTotalCards(cards []CardWithQuantity) int {
	total := 0
	for _, c := range cards {
		total += c.Quantity
	}
	return total
}

// Counts cards by type using normalized string matching.
func countCardTypes(cards []CardWithQuantity, targetType string) int {
	count := 0
	for _, c := range cards {
		if strings.Contains(strings.ToLower(c.Card.Type), strings.ToLower(targetType)) {
			count += c.Quantity
		}
	}
	return count
}

// Tallies monster card attributes.
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

// Computes average ATK and DEF for monster cards, weighted by quantity.
func computeAverageStats(cards []CardWithQuantity) AvgStats {
	var (
		totalATK, totalDEF, monsterCount int
	)

	for _, c := range cards {
		mc := c.Card.MonsterCard
		if mc != nil {
			qty := c.Quantity
			if mc.Atk >= 0 {
				totalATK += mc.Atk * qty
			}
			if mc.Def >= 0 {
				totalDEF += mc.Def * qty
			}
			monsterCount += qty
		}
	}

	if monsterCount == 0 {
		return AvgStats{}
	}

	return AvgStats{
		AvgATK: float64(totalATK) / float64(monsterCount),
		AvgDEF: float64(totalDEF) / float64(monsterCount),
	}
}
