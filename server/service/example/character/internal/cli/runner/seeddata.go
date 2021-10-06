package runner

import (
	"github.com/urfave/cli/v2"

	"gitlab.com/alienspaces/go-boilerplate/server/service/character/internal/harness"
	"gitlab.com/alienspaces/go-boilerplate/server/service/character/internal/record"
)

// LoadSeedData -
func (rnr *Runner) LoadSeedData(c *cli.Context) error {

	rnr.Log.Info("** Load Seed Data **")

	// harness
	config := harness.DataConfig{

		CharacterConfig: []harness.CharacterConfig{
			// Mage - Dark Armoured
			{
				Record: record.Character{
					Name:         "Dark Armoured",
					Strength:     16,
					Dexterity:    12,
					Intelligence: 10,
					Coins:        50,
				},
			},
			// Mage - Red Stripe Druid
			{
				Record: record.Character{
					Name:         "Druid",
					Strength:     14,
					Dexterity:    14,
					Intelligence: 10,
					Coins:        50,
				},
			},
			// Mage - Red Fairy
			{
				Record: record.Character{
					Name:         "Fairy",
					Strength:     10,
					Dexterity:    14,
					Intelligence: 14,
					Coins:        50,
				},
			},
			// Mage - Red Stripe Necromancer
			{
				Record: record.Character{
					Name:         "Necromancer",
					Strength:     14,
					Dexterity:    10,
					Intelligence: 14,
					Coins:        50,
				},
			},
			// Mage - Green Elven
			{
				Record: record.Character{
					Name:         "Elven",
					Strength:     12,
					Dexterity:    14,
					Intelligence: 12,
					Coins:        50,
				},
			},
			// Familliar - Brown Cyclops Bat
			{
				Record: record.Character{
					Name:         "Brown Cyclops Bat",
					Strength:     10,
					Dexterity:    10,
					Intelligence: 10,
				},
			},
			// Familliar - Brown Yeti
			{
				Record: record.Character{
					Name:         "Brown Yeti",
					Strength:     10,
					Dexterity:    10,
					Intelligence: 10,
				},
			},
			// Familliar - Green Tribble
			{
				Record: record.Character{
					Name:         "Green Tribble",
					Strength:     10,
					Dexterity:    10,
					Intelligence: 10,
				},
			},
			// Familliar - Grey Cyclops
			{
				Record: record.Character{
					Name:         "Grey Cyclops",
					Strength:     10,
					Dexterity:    10,
					Intelligence: 10,
				},
			},
			// Familliar - Orange Spotted Tribble
			{
				Record: record.Character{
					Name:         "Orange Spotten Tribble",
					Strength:     10,
					Dexterity:    10,
					Intelligence: 10,
				},
			},
			// Familliar - Purple Bat
			{
				Record: record.Character{
					Name:         "Purple Bat",
					Strength:     10,
					Dexterity:    10,
					Intelligence: 10,
				},
			},
			// Familliar - Purple Minotaur
			{
				Record: record.Character{
					Name:         "Purple Minotaur",
					Strength:     10,
					Dexterity:    10,
					Intelligence: 10,
				},
			},
		},
	}

	h, err := harness.NewTesting(config)
	if err != nil {
		rnr.Log.Warn("Failed new testing harness >%v<", err)
		return err
	}

	// harness commit data
	h.CommitData = true

	err = h.Setup()
	if err != nil {
		rnr.Log.Warn("Failed testing harness setup >%v<", err)
		return err
	}

	rnr.Log.Info("All done")

	return nil
}
