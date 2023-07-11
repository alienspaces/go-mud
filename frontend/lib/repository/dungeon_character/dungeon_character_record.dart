// Package
part of 'dungeon_character_repository.dart';

class DungeonCharacterRecord extends Equatable {
  final String dungeonID;
  final String dungeonName;
  final String dungeonDescription;
  final String locationID;
  final String locationName;
  final String locationDescription;
  final String characterID;
  final String characterName;
  final int characterStrength;
  final int characterDexterity;
  final int characterIntelligence;
  final int characterCurrentStrength;
  final int characterCurrentDexterity;
  final int characterCurrentIntelligence;
  final int characterHealth;
  final int characterFatigue;
  final int characterCurrentHealth;
  final int characterCurrentFatigue;
  final int characterCoins;
  final int characterExperiencePoints;
  final int characterAttributePoints;

  const DungeonCharacterRecord({
    required this.dungeonID,
    required this.dungeonName,
    required this.dungeonDescription,
    required this.locationID,
    required this.locationName,
    required this.locationDescription,
    required this.characterID,
    required this.characterName,
    required this.characterStrength,
    required this.characterDexterity,
    required this.characterIntelligence,
    required this.characterCurrentStrength,
    required this.characterCurrentDexterity,
    required this.characterCurrentIntelligence,
    required this.characterHealth,
    required this.characterFatigue,
    required this.characterCurrentHealth,
    required this.characterCurrentFatigue,
    required this.characterCoins,
    required this.characterExperiencePoints,
    required this.characterAttributePoints,
  });

  DungeonCharacterRecord.fromJson(Map<String, dynamic> json)
      : dungeonID = json['dungeon'] != null ? json['dungeon']['id'] : null,
        dungeonName = json['dungeon'] != null ? json['dungeon']['name'] : null,
        dungeonDescription =
            json['dungeon'] != null ? json['dungeon']['description'] : null,
        locationID = json['location'] != null ? json['location']['id'] : null,
        locationName =
            json['location'] != null ? json['location']['name'] : null,
        locationDescription =
            json['location'] != null ? json['location']['description'] : null,
        characterID = json['id'],
        characterName = json['name'],
        characterStrength = json['strength'],
        characterDexterity = json['dexterity'],
        characterIntelligence = json['intelligence'],
        characterCurrentStrength = json['current_strength'],
        characterCurrentDexterity = json['current_dexterity'],
        characterCurrentIntelligence = json['current_intelligence'],
        characterHealth = json['health'],
        characterFatigue = json['fatigue'],
        characterCurrentHealth = json['current_health'],
        characterCurrentFatigue = json['current_fatigue'],
        characterCoins = json['coins'],
        characterExperiencePoints = json['experience_points'],
        characterAttributePoints = json['attribute_points'];

  Map<String, dynamic> toJson() => {
        'dungeon': {
          'id': dungeonID,
          'name': dungeonName,
          'description': dungeonDescription,
        },
        'location': {
          'id': locationID,
          'name': locationName,
          'description': locationDescription,
        },
        'id': characterID,
        'name': characterName,
        'strength': characterStrength,
        'dexterity': characterDexterity,
        'intelligence': characterIntelligence,
        'current_strength': characterCurrentStrength,
        'current_dexterity': characterCurrentDexterity,
        'current_intelligence': characterCurrentIntelligence,
        'health': characterHealth,
        'fatigue': characterFatigue,
        'current_health': characterCurrentHealth,
        'current_fatigue': characterCurrentFatigue,
        'coins': characterCoins,
        'experience_points': characterExperiencePoints,
        'attribute_points': characterAttributePoints,
      };

  @override
  List<Object?> get props => [
        dungeonID,
        dungeonName,
        dungeonDescription,
        locationID,
        locationName,
        locationDescription,
        characterID,
        characterName,
        characterStrength,
        characterDexterity,
        characterIntelligence,
        characterCurrentStrength,
        characterCurrentDexterity,
        characterCurrentIntelligence,
        characterHealth,
        characterFatigue,
        characterCurrentHealth,
        characterCurrentFatigue,
        characterCoins,
        characterExperiencePoints,
        characterAttributePoints,
      ];
}
