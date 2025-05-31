package services

import (
	"testing"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"github.com/stretchr/testify/assert"
)

func TestCountCardTypes(t *testing.T) {
	cards := []CardWithQuantity{
		{Card: models.Card{Type: "Monster"}, Quantity: 3},
		{Card: models.Card{Type: "Spell Card"}, Quantity: 2},
		{Card: models.Card{Type: "Trap"}, Quantity: 1},
	}

	monsterCount := countCardTypes(cards, "monster")
	spellCount := countCardTypes(cards, "spell")
	trapCount := countCardTypes(cards, "trap")

	assert.Equal(t, 3, monsterCount)
	assert.Equal(t, 2, spellCount)
	assert.Equal(t, 1, trapCount)
}

func TestCountMonsterAttributes(t *testing.T) {
	cards := []CardWithQuantity{
		{
			Quantity: 2,
			Card: models.Card{
				Type:        "Monster",
				MonsterCard: &models.MonsterCard{Attribute: "DARK"},
			},
		},
		{
			Quantity: 1,
			Card: models.Card{
				Type:        "Monster",
				MonsterCard: &models.MonsterCard{Attribute: "LIGHT"},
			},
		},
	}

	attrs := countMonsterAttributes(cards)

	assert.Equal(t, 2, attrs["DARK"])
	assert.Equal(t, 1, attrs["LIGHT"])
}

func TestComputeAverageStats(t *testing.T) {
	cards := []CardWithQuantity{
		{
			Quantity: 2,
			Card: models.Card{
				Type:        "Monster",
				MonsterCard: &models.MonsterCard{Atk: 2000, Def: 1500},
			},
		},
		{
			Quantity: 1,
			Card: models.Card{
				Type:        "Monster",
				MonsterCard: &models.MonsterCard{Atk: 1000, Def: 1200},
			},
		},
	}

	avg := computeAverageStats(cards)

	assert.InEpsilon(t, 1666.6, avg.AvgATK, 0.1)
	assert.InEpsilon(t, 1400.0, avg.AvgDEF, 0.1)
}
