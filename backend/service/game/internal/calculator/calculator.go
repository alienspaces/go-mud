package calculator

import (
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

func CalculateCharacterHealth(rec *record.Character) (*record.Character, error) {

	rec.Health = rec.Strength + rec.Dexterity*10

	return rec, nil
}

func CalculateCharacterFatigue(rec *record.Character) (*record.Character, error) {

	rec.Fatigue = rec.Strength + rec.Intelligence*10

	return rec, nil
}

func CalculateMonsterHealth(rec *record.Monster) (*record.Monster, error) {

	rec.Health = rec.Strength + rec.Dexterity*10

	return rec, nil
}

func CalculateMonsterFatigue(rec *record.Monster) (*record.Monster, error) {

	rec.Fatigue = rec.Strength + rec.Intelligence*10

	return rec, nil
}

func CalculateCharacterDamage(rec *record.CharacterInstanceView, set *record.LocationInstanceViewRecordSet) (int, error) {

	dmg := rec.Strength / 2

	return dmg, nil
}

func CalculateMonsterDamage(rec *record.MonsterInstanceView, set *record.LocationInstanceViewRecordSet) (int, error) {

	dmg := rec.Strength / 2

	return dmg, nil
}
