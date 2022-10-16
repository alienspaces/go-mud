// Package
part of 'dungeon_character_repository.dart';

class DungeonCharacterRecord extends Equatable {
  final String id;
  final String dungeonID;
  final String dungeonName;
  final String dungeonDescription;
  final String locationID;
  final String locationName;
  final String locationDescription;
  final String name;
  final int strength;
  final int dexterity;
  final int intelligence;
  final int currentStrength;
  final int currentDexterity;
  final int currentIntelligence;
  final int health;
  final int fatigue;
  final int currentHealth;
  final int currentFatigue;
  final int coins;
  final int experiencePoints;
  final int attributePoints;

  const DungeonCharacterRecord({
    required this.id,
    required this.dungeonID,
    required this.dungeonName,
    required this.dungeonDescription,
    required this.locationID,
    required this.locationName,
    required this.locationDescription,
    required this.name,
    required this.strength,
    required this.dexterity,
    required this.intelligence,
    required this.currentStrength,
    required this.currentDexterity,
    required this.currentIntelligence,
    required this.health,
    required this.fatigue,
    required this.currentHealth,
    required this.currentFatigue,
    required this.coins,
    required this.experiencePoints,
    required this.attributePoints,
  });

  DungeonCharacterRecord.fromJson(Map<String, dynamic> json)
      : id = json['id'],
        dungeonID = json['dungeon_id'],
        dungeonName = json['dungeon_name'],
        dungeonDescription = json['dungeon_description'],
        locationID = json['location_id'],
        locationName = json['location_name'],
        locationDescription = json['location_description'],
        name = json['name'],
        strength = json['strength'],
        dexterity = json['dexterity'],
        intelligence = json['intelligence'],
        currentStrength = json['current_strength'],
        currentDexterity = json['current_dexterity'],
        currentIntelligence = json['current_intelligence'],
        health = json['health'],
        fatigue = json['fatigue'],
        currentHealth = json['current_health'],
        currentFatigue = json['current_fatigue'],
        coins = json['coins'],
        experiencePoints = json['experiencePoints'],
        attributePoints = json['attributePoints'];

  Map<String, dynamic> toJson() => {
        'id': id,
        'name': name,
        'strength': strength,
        'dexterity': dexterity,
        'intelligence': intelligence,
        'current_strength': currentStrength,
        'current_dexterity': currentDexterity,
        'current_intelligence': currentIntelligence,
        'health': health,
        'fatigue': fatigue,
        'current_health': currentHealth,
        'current_fatigue': currentFatigue,
        'coins': coins,
        'experiencePoints': experiencePoints,
        'attributePoints': attributePoints,
      };

  @override
  List<Object?> get props => [id, name];
}
