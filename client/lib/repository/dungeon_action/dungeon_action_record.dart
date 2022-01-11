// Package
part of 'dungeon_action_repository.dart';

class CreateDungeonActionRecord extends Equatable {
  final String sentence;

  const CreateDungeonActionRecord({
    required this.sentence,
  });

  CreateDungeonActionRecord.fromJson(Map<String, dynamic> json)
      : sentence = json['sentence'];

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
  final CharacterDetailedData? character;
  final MonsterDetailedData? monster;
  final ObjectDetailedData? equippedObject;
  final ObjectDetailedData? stashedObject;
  final ObjectDetailedData? targetObject;
  final CharacterDetailedData? targetCharacter;
  final MonsterDetailedData? targetMonster;
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

    // Location objects
    List<dynamic>? locationObjects = location['objects'];
    List<ObjectData>? locationObjectData;
    if (locationObjects != null) {
      locationObjectData =
          locationObjects.map((e) => ObjectData.fromJson(e)).toList();
    }

    // Location characters
    List<dynamic>? locationCharacters = location['characters'];
    List<CharacterData>? locationCharacterData;
    if (locationCharacters != null) {
      locationCharacterData =
          locationCharacters.map((e) => CharacterData.fromJson(e)).toList();
    }

// Location monsters
    List<dynamic>? locationMonsters = location['monsters'];
    List<MonsterData>? locationMonsterData;
    if (locationMonsters != null) {
      locationMonsterData =
          locationMonsters.map((e) => MonsterData.fromJson(e)).toList();
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

    MonsterDetailedData? monsterData;
    if (monster != null) {
      monsterData = MonsterDetailedData.fromJson(monster);
    }

    CharacterDetailedData? characterData;
    if (character != null) {
      characterData = CharacterDetailedData.fromJson(character);
    }

    // Equipped object
    Map<String, dynamic>? equippedObject = json['equipped_object'];
    ObjectDetailedData? equippedObjectData;
    if (equippedObject != null) {
      equippedObjectData = ObjectDetailedData.fromJson(equippedObject);
    }

    // Stashed object
    Map<String, dynamic>? stashedObject = json['stashed_object'];
    ObjectDetailedData? stashedObjectData;
    if (stashedObject != null) {
      stashedObjectData = ObjectDetailedData.fromJson(stashedObject);
    }

    // Target object
    Map<String, dynamic>? targetObject = json['target_object'];
    ObjectDetailedData? targetObjectData;
    if (targetObject != null) {
      targetObjectData = ObjectDetailedData.fromJson(targetObject);
    }

    // Target character
    Map<String, dynamic>? targetCharacter = json['target_character'];
    CharacterDetailedData? targetCharacterData;
    if (targetCharacter != null) {
      targetCharacterData = CharacterDetailedData.fromJson(targetCharacter);
    }

    // Target monster
    Map<String, dynamic>? targetMonster = json['target_monster'];
    MonsterDetailedData? targetMonsterData;
    if (targetMonster != null) {
      targetMonsterData = MonsterDetailedData.fromJson(targetMonster);
    }

    // Target location
    LocationData? targetLocationData;

    Map<String, dynamic>? targetLocation = json['target_location'];
    if (targetLocation != null) {
      List<dynamic>? locationObjects = targetLocation['objects'];
      List<ObjectData>? locationObjectData;
      if (locationObjects != null) {
        locationObjectData =
            locationObjects.map((e) => ObjectData.fromJson(e)).toList();
      }

      List<dynamic>? locationCharacters = targetLocation['characters'];
      List<CharacterData>? locationCharacterData;
      if (locationCharacters != null) {
        locationCharacterData =
            locationCharacters.map((e) => CharacterData.fromJson(e)).toList();
      }

      List<dynamic>? locationMonsters = targetLocation['monsters'];
      List<MonsterData>? locationMonsterData;
      if (locationMonsters != null) {
        locationMonsterData =
            locationMonsters.map((e) => MonsterData.fromJson(e)).toList();
      }

      List<dynamic> directions = targetLocation['directions'];

      targetLocationData = LocationData(
        name: targetLocation['name'],
        description: targetLocation['description'],
        direction: targetLocation['direction'],
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

class ObjectDetailedData {
  final String name;
  final String description;
  final bool isStashed;
  final bool isEquipped;
  ObjectDetailedData(
      {required this.name,
      required this.description,
      required this.isStashed,
      required this.isEquipped});

  factory ObjectDetailedData.fromJson(Map<String, dynamic> json) {
    return ObjectDetailedData(
      name: json['name'],
      description: json['description'],
      isStashed: json['is_stashed'],
      isEquipped: json['is_equipped'],
    );
  }
}

class CharacterData {
  final String name;

  CharacterData({
    required this.name,
  });

  factory CharacterData.fromJson(Map<String, dynamic> json) {
    return CharacterData(
      name: json['name'],
    );
  }
}

class CharacterDetailedData {
  final String name;
  final int strength;
  final int dexterity;
  final int intelligence;
  final int currentStrength;
  final int currentDexterity;
  final int currentIntelligence;
  final int health;
  final int fatigue;

  CharacterDetailedData({
    required this.name,
    required this.strength,
    required this.dexterity,
    required this.intelligence,
    required this.currentStrength,
    required this.currentDexterity,
    required this.currentIntelligence,
    required this.health,
    required this.fatigue,
  });

  factory CharacterDetailedData.fromJson(Map<String, dynamic> json) {
    return CharacterDetailedData(
      name: json['name'],
      strength: json['strength'],
      dexterity: json['dexterity'],
      intelligence: json['intelligence'],
      currentStrength: json['current_strength'],
      currentDexterity: json['current_dexterity'],
      currentIntelligence: json['current_intelligence'],
      health: json['health'],
      fatigue: json['fatigue'],
    );
  }
}

class MonsterData {
  final String name;

  MonsterData({
    required this.name,
  });

  factory MonsterData.fromJson(Map<String, dynamic> json) {
    return MonsterData(
      name: json['name'],
    );
  }
}

class MonsterDetailedData {
  final String name;
  final int strength;
  final int dexterity;
  final int intelligence;
  final int currentStrength;
  final int currentDexterity;
  final int currentIntelligence;
  final int health;
  final int fatigue;

  MonsterDetailedData({
    required this.name,
    required this.strength,
    required this.dexterity,
    required this.intelligence,
    required this.currentStrength,
    required this.currentDexterity,
    required this.currentIntelligence,
    required this.health,
    required this.fatigue,
  });

  factory MonsterDetailedData.fromJson(Map<String, dynamic> json) {
    return MonsterDetailedData(
      name: json['name'],
      strength: json['strength'],
      dexterity: json['dexterity'],
      intelligence: json['intelligence'],
      currentStrength: json['current_strength'],
      currentDexterity: json['current_dexterity'],
      currentIntelligence: json['current_intelligence'],
      health: json['health'],
      fatigue: json['fatigue'],
    );
  }
}
