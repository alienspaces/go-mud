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
      : dungeonID = json['dungeon_id'],
        dungeonName = json['dungeon_name'],
        dungeonDescription = json['dungeon_description'],
        locationID = json['location_id'],
        locationName = json['location_name'],
        locationDescription = json['location_description'],
        characterID = json['character_id'],
        characterName = json['character_name'],
        characterStrength = json['character_strength'],
        characterDexterity = json['character_dexterity'],
        characterIntelligence = json['character_intelligence'],
        characterCurrentStrength = json['character_current_strength'],
        characterCurrentDexterity = json['character_current_dexterity'],
        characterCurrentIntelligence = json['character_current_intelligence'],
        characterHealth = json['character_health'],
        characterFatigue = json['character_fatigue'],
        characterCurrentHealth = json['character_current_health'],
        characterCurrentFatigue = json['character_current_fatigue'],
        characterCoins = json['character_coins'] ?? 0,
        characterExperiencePoints = json['character_experience_points'] ?? 0,
        characterAttributePoints = json['character_attribute_points'] ?? 0;

  Map<String, dynamic> toJson() => {
        'dungeon_id': dungeonID,
        'dungeon_name': dungeonName,
        'dungeon_description': dungeonDescription,
        'location_id': locationID,
        'location_name': locationName,
        'location_description': locationDescription,
        'character_id': characterID,
        'character_name': characterName,
        'character_strength': characterStrength,
        'character_dexterity': characterDexterity,
        'character_intelligence': characterIntelligence,
        'character_current_strength': characterCurrentStrength,
        'character_current_dexterity': characterCurrentDexterity,
        'character_current_intelligence': characterCurrentIntelligence,
        'character_health': characterHealth,
        'character_fatigue': characterFatigue,
        'character_current_health': characterCurrentHealth,
        'character_current_fatigue': characterCurrentFatigue,
        'character_coins': characterCoins,
        'character_experience_points': characterExperiencePoints,
        'character_attribute_points': characterAttributePoints,
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
