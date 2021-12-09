// Package
part of 'dungeon_action_repository.dart';

class CreateDungeonActionRecord extends Equatable {
  final String sentence;

  const CreateDungeonActionRecord({
    required this.sentence,
  });

  CreateDungeonActionRecord.fromJson(Map<String, dynamic> json) : sentence = json['sentence'];

  Map<String, dynamic> toJson() => {
        'sentence': sentence,
      };

  @override
  List<Object> get props => [sentence];
}

class DungeonActionRecord extends Equatable {
  final String id;
  final String command;
  final String commandResult;
  final LocationData location;
  final CharacterData? character;
  final MonsterData? monster;
  final ObjectData? equippedObject;
  final ObjectData? stashedObject;
  final ObjectData? targetObject;
  final CharacterData? targetCharacter;
  final MonsterData? targetMonster;
  final LocationData? targetLocation;

  const DungeonActionRecord({
    required this.id,
    required this.command,
    required this.commandResult,
    required this.location,
    required this.character,
    required this.monster,
    required this.equippedObject,
    required this.stashedObject,
    required this.targetObject,
    required this.targetCharacter,
    required this.targetMonster,
    required this.targetLocation,
  });

  factory DungeonActionRecord.fromJson(Map<String, dynamic> json) {
    // Location
    Map<String, dynamic>? location = json['location'];
    if (location == null) {
      throw const FormatException('Missing "location" from JSON');
    }

    List<dynamic>? locationObjects = location['objects'];
    List<ObjectData>? locationObjectData;
    if (locationObjects != null) {
      locationObjectData = locationObjects.map((e) => ObjectData.fromJson(e)).toList();
    }

    List<dynamic>? locationCharacters = location['characters'];
    List<CharacterData>? locationCharacterData;
    if (locationCharacters != null) {
      locationCharacterData = locationCharacters.map((e) => CharacterData.fromJson(e)).toList();
    }

    List<dynamic>? locationMonsters = location['monsters'];
    List<MonsterData>? locationMonsterData;
    if (locationMonsters != null) {
      locationMonsterData = locationMonsters.map((e) => MonsterData.fromJson(e)).toList();
    }

    List<dynamic> directions = location['directions'];

    var locationData = LocationData(
      name: location['name'],
      description: location['description'],
      direction: location['direction'],
      directions: directions.map((e) => e.toString()).toList(),
      characters: locationCharacterData,
      monsters: locationMonsterData,
      objects: locationObjectData,
    );

    // Character or Monster
    Map<String, dynamic>? character = json['character'];
    Map<String, dynamic>? monster = json['monster'];
    if (character == null && monster == null) {
      throw const FormatException('Missing "character" or "monster" from JSON');
    }

    MonsterData? monsterData;
    if (monster != null) {
      monsterData = MonsterData.fromJson(monster);
    }

    CharacterData? characterData;
    if (character != null) {
      characterData = CharacterData.fromJson(character);
    }

    // Equipped object
    Map<String, dynamic>? equippedObject = json['equipped_object'];
    ObjectData? equippedObjectData;
    if (equippedObject != null) {
      equippedObjectData = ObjectData.fromJson(equippedObject);
    }

    // Stashed object
    Map<String, dynamic>? stashedObject = json['stashed_object'];
    ObjectData? stashedObjectData;
    if (stashedObject != null) {
      stashedObjectData = ObjectData.fromJson(stashedObject);
    }

    // Target object
    Map<String, dynamic>? targetObject = json['target_object'];
    ObjectData? targetObjectData;
    if (targetObject != null) {
      targetObjectData = ObjectData.fromJson(targetObject);
    }

    // Target character
    Map<String, dynamic>? targetCharacter = json['target_character'];
    CharacterData? targetCharacterData;
    if (targetCharacter != null) {
      targetCharacterData = CharacterData.fromJson(targetCharacter);
    }

    // Target monster
    Map<String, dynamic>? targetMonster = json['target_monster'];
    MonsterData? targetMonsterData;
    if (targetMonster != null) {
      targetMonsterData = MonsterData.fromJson(targetMonster);
    }

    // Target location
    LocationData? targetLocationData;

    Map<String, dynamic>? targetLocation = json['target_location'];
    if (targetLocation != null) {
      List<dynamic>? locationObjects = targetLocation['objects'];
      List<ObjectData>? locationObjectData;
      if (locationObjects != null) {
        locationObjectData = locationObjects.map((e) => ObjectData.fromJson(e)).toList();
      }

      List<dynamic>? locationCharacters = targetLocation['characters'];
      List<CharacterData>? locationCharacterData;
      if (locationCharacters != null) {
        locationCharacterData = locationCharacters.map((e) => CharacterData.fromJson(e)).toList();
      }

      List<dynamic>? locationMonsters = targetLocation['monsters'];
      List<MonsterData>? locationMonsterData;
      if (locationMonsters != null) {
        locationMonsterData = locationMonsters.map((e) => MonsterData.fromJson(e)).toList();
      }

      List<dynamic> directions = targetLocation['directions'];

      targetLocationData = LocationData(
        name: location['name'],
        description: location['description'],
        direction: location['direction'],
        directions: directions.map((e) => e.toString()).toList(),
        characters: locationCharacterData,
        monsters: locationMonsterData,
        objects: locationObjectData,
      );
    }

    return DungeonActionRecord(
      id: json['id'],
      command: json['command'],
      commandResult: json['command_result'],
      location: locationData,
      character: characterData,
      monster: monsterData,
      equippedObject: equippedObjectData,
      stashedObject: stashedObjectData,
      targetObject: targetObjectData,
      targetCharacter: targetCharacterData,
      targetMonster: targetMonsterData,
      targetLocation: targetLocationData,
    );
  }

  @override
  List<Object?> get props => [
        id,
        command,
        commandResult,
        location,
        character,
        monster,
        equippedObject,
        stashedObject,
        targetObject,
        targetCharacter,
        targetMonster,
        targetLocation,
      ];
}

class LocationData {
  final String name;
  final String description;
  final String? direction;
  final List<String> directions;
  final List<CharacterData>? characters;
  final List<MonsterData>? monsters;
  final List<ObjectData>? objects;
  LocationData({
    required this.name,
    required this.description,
    this.direction,
    required this.directions,
    this.characters,
    this.monsters,
    this.objects,
  });
}

class ObjectData {
  final String name;
  ObjectData({required this.name});

  factory ObjectData.fromJson(Map<String, dynamic> json) {
    return ObjectData(name: json['name']);
  }
}

class CharacterData {
  final String name;
  CharacterData({required this.name});

  factory CharacterData.fromJson(Map<String, dynamic> json) {
    return CharacterData(name: json['name']);
  }
}

class MonsterData {
  final String name;
  MonsterData({required this.name});

  factory MonsterData.fromJson(Map<String, dynamic> json) {
    return MonsterData(name: json['name']);
  }
}
