package record

type ActionRecordSet struct {
	ActionRec *Action
	// Memories the action
	ActionMemoryRecs []*ActionMemory
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
