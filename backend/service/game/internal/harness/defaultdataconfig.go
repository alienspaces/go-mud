package harness

import (
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

const (
	CharacterNameBarricade string = "Barricade"
	CharacterNameLegislate string = "Legislate"
	CharacterNameBolster   string = "Bolster"
)

const (
	MonsterNameGrumpyDwarf string = "Grumpy Dwarf"
	MonsterNameAngryGoblin string = "Angry Goblin"
)

const (
	ObjectNameRustedSword        string = "Rusted Sword"
	ObjectNameRustedHelmet       string = "Rusted Helmet"
	ObjectNameSilverKey          string = "Silver Key"
	ObjectNameDullBronzeRing     string = "Dull Bronze Ring"
	ObjectNameBloodStainedPouch  string = "Blood Stained Pouch"
	ObjectNameBoneDagger         string = "Bone Dagger"
	ObjectNameVialOfOgreBlood    string = "Vial Of Ogre Blood"
	ObjectNameStoneMace          string = "Stone Mace"
	ObjectNameChippedHammer      string = "Chipped Hammer"
	ObjectNameChippedBreastplate string = "Chipped Breastplate"
)

const (
	DungeonNameCave string = "Cave"
)

const (
	LocationNameCaveEntrance     string = "Cave Entrance"
	LocationNameCaveTunnel       string = "Cave Tunnel"
	LocationNameCaveRoom         string = "Cave Room"
	LocationNameNarrowTunnel     string = "Narrow Tunnel"
	LocationNameDarkNarrowTunnel string = "Dark Narrow Tunnel"
	LocationNameDarkRoom         string = "Dark Room"
)

var DefaultDataConfig = DataConfig{
	ObjectConfig: []ObjectConfig{
		{
			Record: record.Object{
				Name:                ObjectNameRustedSword,
				Description:         "A rusted sword.",
				DescriptionDetailed: "A rusted sword with a chipped blade and a worn leather handle.",
			},
		},
		{
			Record: record.Object{
				Name:                ObjectNameRustedHelmet,
				Description:         "A rusted helmet.",
				DescriptionDetailed: "A rusted helmet pitted with dents.",
			},
		},
		{
			Record: record.Object{
				Name:                ObjectNameSilverKey,
				Description:         "A silver key.",
				DescriptionDetailed: "A silver key with fine runes in a language you do not understand engraved along the edge.",
			},
		},
		{
			Record: record.Object{
				Name:                ObjectNameDullBronzeRing,
				Description:         "A dull bronze ring.",
				DescriptionDetailed: "A dull bronze ring.",
			},
		},
		{
			Record: record.Object{
				Name:                ObjectNameBloodStainedPouch,
				Description:         "A blood stained pouch.",
				DescriptionDetailed: "A blood stained pouch.",
			},
		},
		{
			Record: record.Object{
				Name:                ObjectNameBoneDagger,
				Description:         "A bone dagger.",
				DescriptionDetailed: "A bone dagger.",
			},
		},
		{
			Record: record.Object{
				Name:                ObjectNameVialOfOgreBlood,
				Description:         "A large vial of ogre blood.",
				DescriptionDetailed: "A large vial of ogre blood.",
			},
		},
		{
			Record: record.Object{
				Name:                ObjectNameStoneMace,
				Description:         "A stone mace.",
				DescriptionDetailed: "A stone mace.",
			},
		},
		{
			Record: record.Object{
				Name:                ObjectNameChippedHammer,
				Description:         "A chipped hammer.",
				DescriptionDetailed: "A chipped hammer.",
			},
		},
		{
			Record: record.Object{
				Name:                ObjectNameChippedBreastplate,
				Description:         "A chipped breastplate.",
				DescriptionDetailed: "A chipped breastplate.",
			},
		},
	},
	MonsterConfig: []MonsterConfig{
		{
			Record: record.Monster{
				Name:        MonsterNameGrumpyDwarf,
				Description: "A particularly grumpy specimen of a common dwarf.",
			},
			MonsterObjectConfig: []MonsterObjectConfig{
				{
					Record: record.MonsterObject{
						IsEquipped: true,
					},
					ObjectName: ObjectNameBoneDagger,
				},
				{
					Record: record.MonsterObject{
						IsStashed: true,
					},
					ObjectName: ObjectNameVialOfOgreBlood,
				},
			},
		},
		{
			Record: record.Monster{
				Name:        MonsterNameAngryGoblin,
				Description: "A particularly angry specimen of a common goblin.",
			},
			MonsterObjectConfig: []MonsterObjectConfig{},
		},
	},
	CharacterConfig: []CharacterConfig{
		{
			Record: record.Character{
				Name: CharacterNameBarricade,
			},
			CharacterObjectConfig: []CharacterObjectConfig{
				{
					Record: record.CharacterObject{
						IsEquipped: true,
					},
					ObjectName: ObjectNameDullBronzeRing,
				},
				{
					Record: record.CharacterObject{
						IsStashed: true,
					},
					ObjectName: ObjectNameBloodStainedPouch,
				},
			},
		},
		{
			Record: record.Character{
				Name: CharacterNameLegislate,
			},
			CharacterObjectConfig: []CharacterObjectConfig{
				{
					Record: record.CharacterObject{
						IsEquipped: true,
					},
					ObjectName: ObjectNameStoneMace,
				},
			},
		},
		{
			Record: record.Character{
				Name: CharacterNameBolster,
			},
			CharacterObjectConfig: []CharacterObjectConfig{
				{
					Record: record.CharacterObject{
						IsEquipped: true,
					},
					ObjectName: ObjectNameChippedHammer,
				},
				{
					Record: record.CharacterObject{
						IsEquipped: true,
					},
					ObjectName: ObjectNameChippedBreastplate,
				},
			},
		},
	},
	DungeonConfig: []DungeonConfig{
		{
			Record: record.Dungeon{
				Name:        DungeonNameCave,
				Description: "A dark and damp stone cave.",
			},
			LocationConfig: []LocationConfig{
				{
					Record: record.Location{
						Name:        LocationNameCaveEntrance,
						Description: "A large cave entrance.",
						IsDefault:   true,
					},
					NorthLocationName: LocationNameCaveTunnel,
					LocationMonsterConfig: []LocationMonsterConfig{
						{
							MonsterName: MonsterNameGrumpyDwarf,
						},
					},
					LocationObjectConfig: []LocationObjectConfig{
						{
							ObjectName: ObjectNameRustedSword,
						},
					},
				},
				{
					Record: record.Location{
						Name:        LocationNameCaveTunnel,
						Description: "A cave tunnel descends into the mountain.",
					},
					NorthLocationName:     LocationNameCaveRoom,
					SouthLocationName:     LocationNameCaveEntrance,
					NorthwestLocationName: LocationNameNarrowTunnel,
					LocationMonsterConfig: []LocationMonsterConfig{
						{
							MonsterName: MonsterNameAngryGoblin,
						},
					},
					LocationObjectConfig: []LocationObjectConfig{
						{
							ObjectName: ObjectNameRustedHelmet,
						},
					},
				},
				{
					Record: record.Location{
						Name:        LocationNameCaveRoom,
						Description: "A large cave room.",
					},
					SouthLocationName:     LocationNameCaveTunnel,
					LocationMonsterConfig: []LocationMonsterConfig{},
					LocationObjectConfig: []LocationObjectConfig{
						{
							ObjectName: ObjectNameSilverKey,
						},
					},
				},
				{
					Record: record.Location{
						Name:        LocationNameNarrowTunnel,
						Description: "A narrow tunnel gradually descending into the darkness.",
					},
					NorthwestLocationName: LocationNameDarkNarrowTunnel,
					SoutheastLocationName: LocationNameCaveTunnel,
					LocationMonsterConfig: []LocationMonsterConfig{},
					LocationObjectConfig:  []LocationObjectConfig{},
				},
				{
					Record: record.Location{
						Name:        LocationNameDarkNarrowTunnel,
						Description: "A dark narrow tunnel.",
					},
					SoutheastLocationName: LocationNameNarrowTunnel,
					DownLocationName:      LocationNameDarkRoom,
					LocationMonsterConfig: []LocationMonsterConfig{},
					LocationObjectConfig:  []LocationObjectConfig{},
				},
				{
					Record: record.Location{
						Name:        LocationNameDarkRoom,
						Description: "A dark room.",
					},
					UpLocationName:        LocationNameDarkNarrowTunnel,
					LocationMonsterConfig: []LocationMonsterConfig{},
					LocationObjectConfig:  []LocationObjectConfig{},
				},
			},
			DungeonInstanceConfig: []DungeonInstanceConfig{
				{
					CharacterInstanceConfig: []CharacterInstanceConfig{
						{
							Name: CharacterNameBarricade,
						},
						{
							Name: CharacterNameLegislate,
						},
					},
					TurnConfig: []TurnConfig{
						{
							ActionConfig: []ActionConfig{
								{
									CharacterName: CharacterNameBarricade,
									Command:       "look",
								},
								{
									MonsterName: MonsterNameGrumpyDwarf,
									Command:     "look",
								},
							},
						},
					},
				},
			},
		},
	},
}
