// Package
part of 'character_repository.dart';

class CreateCharacterRecord extends Equatable {
  final String name;
  final int strength;
  final int dexterity;
  final int intelligence;

  const CreateCharacterRecord({
    required this.name,
    required this.strength,
    required this.dexterity,
    required this.intelligence,
  });

  CreateCharacterRecord.fromJson(Map<String, dynamic> json)
      : name = json['name'],
        strength = json['strength'],
        dexterity = json['dexterity'],
        intelligence = json['intelligence'];

  Map<String, dynamic> toJson() => {
        'name': name,
        'strength': strength,
        'dexterity': dexterity,
        'intelligence': intelligence,
      };

  @override
  List<Object?> get props => [name, strength, dexterity, intelligence];
}

class CharacterRecord extends Equatable {
  final String id;
  final String name;
  final int strength;
  final int dexterity;
  final int intelligence;
  final int? health;
  final int? fatigue;
  final int? coins;
  final int? experiencePoints;
  final int? attributePoints;

  const CharacterRecord({
    required this.id,
    required this.name,
    required this.strength,
    required this.dexterity,
    required this.intelligence,
    this.health,
    this.fatigue,
    this.coins,
    this.experiencePoints,
    this.attributePoints,
  });

  CharacterRecord.fromJson(Map<String, dynamic> json)
      : id = json['id'],
        name = json['name'],
        strength = json['strength'],
        dexterity = json['dexterity'],
        intelligence = json['intelligence'],
        health = json['health'],
        fatigue = json['fatigue'],
        coins = json['coins'],
        experiencePoints = json['experiencePoints'],
        attributePoints = json['attributePoints'];

  Map<String, dynamic> toJson() => {
        'id': id,
        'name': name,
        'strength': strength,
        'dexterity': dexterity,
        'intelligence': intelligence,
        'health': health,
        'fatigue': fatigue,
        'coins': coins,
        'experiencePoints': experiencePoints,
        'attributePoints': attributePoints,
      };

  @override
  List<Object?> get props => [id, name];
}
