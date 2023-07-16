package record

type ActionRecordSet struct {
	// The current action record
	ActionRec *Action
	// The character performing the action
	ActionCharacterRec *ActionCharacter
	// The stashed and equipped objects of the character that is performing the action
	ActionCharacterObjectRecs []*ActionCharacterObject
	// The monster performing the action
	ActionMonsterRec *ActionMonster
	// The stashed and equipped objects of the monster performing the action
	ActionMonsterObjectRecs []*ActionMonsterObject
	// The current location of the character or monster performing the action
	CurrentLocation *ActionLocationRecordSet
	// The object that was equipped as a result of an action
	EquippedActionObjectRec *ActionObject
	// The object that was stashed as a result of an action
	StashedActionObjectRec *ActionObject
	// The object that was dropped as a result of an action
	DroppedActionObjectRec *ActionObject
	// The object that the action is being performed on
	TargetActionObjectRec *ActionObject
	// The character the action is being performed on
	TargetActionCharacterRec *ActionCharacter
	// The equipped objects of the character the action is being performed on
	TargetActionCharacterObjectRecs []*ActionCharacterObject
	// The monster the action is being performed on
	TargetActionMonsterRec *ActionMonster
	// The equipped objects of the monster the action is being performed on
	TargetActionMonsterObjectRecs []*ActionMonsterObject
	// The location where the action is being performed
	TargetLocation *ActionLocationRecordSet
}

type ActionLocationRecordSet struct {
	LocationInstanceViewRec *LocationInstanceView
	ActionCharacterRecs     []*ActionCharacter
	ActionMonsterRecs       []*ActionMonster
	ActionObjectRecs        []*ActionObject
}

type LocationInstanceViewRecordSet struct {
	LocationInstanceViewRec   *LocationInstanceView
	CharacterInstanceViewRecs []*CharacterInstanceView
	MonsterInstanceViewRecs   []*MonsterInstanceView
	ObjectInstanceViewRecs    []*ObjectInstanceView
	LocationInstanceViewRecs  []*LocationInstanceView
}

func (l *LocationInstanceViewRecordSet) LocationDirections() []string {
	d := []string{}

	if l.LocationInstanceViewRec == nil {
		return d
	}

	li := l.LocationInstanceViewRec
	if li.NorthLocationInstanceID.Valid {
		d = append(d, "north")
	}
	if li.NortheastLocationInstanceID.Valid {
		d = append(d, "northeast")
	}
	if li.EastLocationInstanceID.Valid {
		d = append(d, "east")
	}
	if li.SoutheastLocationInstanceID.Valid {
		d = append(d, "southeast")
	}
	if li.SouthLocationInstanceID.Valid {
		d = append(d, "south")
	}
	if li.SouthwestLocationInstanceID.Valid {
		d = append(d, "southwest")
	}
	if li.WestLocationInstanceID.Valid {
		d = append(d, "west")
	}
	if li.NorthwestLocationInstanceID.Valid {
		d = append(d, "northwest")
	}
	if li.UpLocationInstanceID.Valid {
		d = append(d, "up")
	}
	if li.DownLocationInstanceID.Valid {
		d = append(d, "down")
	}
	return d
}
