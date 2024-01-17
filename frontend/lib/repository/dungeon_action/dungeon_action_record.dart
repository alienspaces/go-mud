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
  final String actionID;
  final String actionCommand;
  final String actionNarrative;
  final int actionTurnNumber;
  final int actionSerialNumber;
  final LocationData actionLocation;
  final CharacterDetailedData? actionCharacter;
  final MonsterDetailedData? actionMonster;
  final ObjectDetailedData? actionEquippedObject;
  final ObjectDetailedData? actionStashedObject;
  final ObjectDetailedData? actionDroppedObject;
  final ObjectDetailedData? actionTargetObject;
  final CharacterDetailedData? actionTargetCharacter;
  final MonsterDetailedData? actionTargetMonster;
  final LocationData? actionTargetLocation;

  const DungeonActionRecord({
    required this.actionID,
    required this.actionCommand,
    required this.actionNarrative,
    required this.actionTurnNumber,
    required this.actionSerialNumber,
    required this.actionLocation,
    required this.actionCharacter,
    required this.actionMonster,
    required this.actionEquippedObject,
    required this.actionStashedObject,
    required this.actionDroppedObject,
    required this.actionTargetObject,
    required this.actionTargetCharacter,
    required this.actionTargetMonster,
    required this.actionTargetLocation,
  });

  factory DungeonActionRecord.fromJson(Map<String, dynamic> json) {
    // Location
    Map<String, dynamic>? location = json['location'];
    if (location == null) {
      throw const FormatException('Missing "action_location" from JSON');
    }

    var locationData = LocationData.fromJson(location);

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

    // Dropped object
    Map<String, dynamic>? droppedObject = json['dropped_object'];
    ObjectDetailedData? droppedObjectData;
    if (droppedObject != null) {
      droppedObjectData = ObjectDetailedData.fromJson(droppedObject);
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
      targetLocationData = LocationData.fromJson(targetLocation);
    }

    return DungeonActionRecord(
      actionID: json['id'],
      actionCommand: json['command'],
      actionNarrative: json['narrative'],
      actionTurnNumber: json['turn_number'],
      actionSerialNumber: json['serial_number'],
      actionLocation: locationData,
      actionCharacter: characterData,
      actionMonster: monsterData,
      actionEquippedObject: equippedObjectData,
      actionStashedObject: stashedObjectData,
      actionDroppedObject: droppedObjectData,
      actionTargetObject: targetObjectData,
      actionTargetCharacter: targetCharacterData,
      actionTargetMonster: targetMonsterData,
      actionTargetLocation: targetLocationData,
    );
  }

  @override
  List<Object?> get props => [
        actionID,
        actionCommand,
        actionNarrative,
        actionLocation,
        actionCharacter,
        actionMonster,
        actionEquippedObject,
        actionStashedObject,
        actionDroppedObject,
        actionTargetObject,
        actionTargetCharacter,
        actionTargetMonster,
        actionTargetLocation,
      ];
}

class LocationData {
  final String locationName;
  final String locationDescription;
  final String? locationDirection;
  final List<String> locationDirections;
  final List<CharacterData>? locationCharacters;
  final List<MonsterData>? locationMonsters;
  final List<ObjectData>? locationObjects;

  LocationData({
    required this.locationName,
    required this.locationDescription,
    this.locationDirection,
    required this.locationDirections,
    this.locationCharacters,
    this.locationMonsters,
    this.locationObjects,
  });

  factory LocationData.fromJson(Map<String, dynamic> json) {
    List<dynamic> directions = json['directions'];

    // Location objects
    List<dynamic>? locationObjects = json['objects'];
    List<ObjectData>? locationObjectData;
    if (locationObjects != null) {
      locationObjectData =
          locationObjects.map((e) => ObjectData.fromJson(e)).toList();
    }

    // Location characters
    List<dynamic>? locationCharacters = json['characters'];
    List<CharacterData>? locationCharacterData;
    if (locationCharacters != null) {
      locationCharacterData =
          locationCharacters.map((e) => CharacterData.fromJson(e)).toList();
    }

    // Location monsters
    List<dynamic>? locationMonsters = json['monsters'];
    List<MonsterData>? locationMonsterData;
    if (locationMonsters != null) {
      locationMonsterData =
          locationMonsters.map((e) => MonsterData.fromJson(e)).toList();
    }

    return LocationData(
      locationName: json['name'],
      locationDescription: json['description'],
      locationDirection: json['direction'],
      locationDirections: directions.map((e) => e.toString()).toList(),
      locationCharacters: locationCharacterData,
      locationMonsters: locationMonsterData,
      locationObjects: locationObjectData,
    );
  }
}

class ObjectData {
  final String name;
  ObjectData({required this.name});

  factory ObjectData.fromJson(Map<String, dynamic> json) {
    return ObjectData(name: json['name']);
  }
}

class ObjectDetailedData {
  final String objectName;
  final String objectDescription;
  final bool objectIsStashed;
  final bool objectIsEquipped;

  ObjectDetailedData(
      {required this.objectName,
      required this.objectDescription,
      required this.objectIsStashed,
      required this.objectIsEquipped});

  factory ObjectDetailedData.fromJson(Map<String, dynamic> json) {
    return ObjectDetailedData(
      objectName: json['name'],
      objectDescription: json['description'],
      objectIsStashed: json['is_stashed'],
      objectIsEquipped: json['is_equipped'],
    );
  }

  Map<String, dynamic> toJson() => {
        'name': objectName,
        'description': objectDescription,
        'is_stashed': objectIsStashed,
        'is_equipped': objectIsEquipped,
      };
}

class CharacterData {
  final String name;
  final int health;
  final int fatigue;
  final int currentHealth;
  final int currentFatigue;

  CharacterData({
    required this.name,
    required this.health,
    required this.fatigue,
    required this.currentHealth,
    required this.currentFatigue,
  });

  factory CharacterData.fromJson(Map<String, dynamic> json) {
    return CharacterData(
      name: json['name'],
      health: json['health'],
      fatigue: json['fatigue'],
      currentHealth: json['current_health'],
      currentFatigue: json['current_fatigue'],
    );
  }
}

class CharacterDetailedData {
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
  final List<ObjectDetailedData>? characterStashedObjects;
  final List<ObjectDetailedData>? characterEquippedObjects;

  CharacterDetailedData({
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
    required this.characterStashedObjects,
    required this.characterEquippedObjects,
  });

  factory CharacterDetailedData.fromJson(Map<String, dynamic> json) {
    List<dynamic>? equippedObjects = json['equipped_objects'];
    List<ObjectDetailedData>? equippedObjectData;
    if (equippedObjects != null) {
      equippedObjectData =
          equippedObjects.map((e) => ObjectDetailedData.fromJson(e)).toList();
    }

    List<dynamic>? stashedObjects = json['stashed_objects'];
    List<ObjectDetailedData>? stashedObjectData;
    if (stashedObjects != null) {
      stashedObjectData =
          stashedObjects.map((e) => ObjectDetailedData.fromJson(e)).toList();
    }

    return CharacterDetailedData(
      characterName: json['name'],
      characterStrength: json['strength'],
      characterDexterity: json['dexterity'],
      characterIntelligence: json['intelligence'],
      characterCurrentStrength: json['current_strength'],
      characterCurrentDexterity: json['current_dexterity'],
      characterCurrentIntelligence: json['current_intelligence'],
      characterHealth: json['health'],
      characterFatigue: json['fatigue'],
      characterCurrentHealth: json['current_health'],
      characterCurrentFatigue: json['current_fatigue'],
      characterStashedObjects: stashedObjectData,
      characterEquippedObjects: equippedObjectData,
    );
  }

  Map<String, dynamic> toJson() => {
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
      };
}

class MonsterData {
  final String name;
  final int health;
  final int fatigue;
  final int currentHealth;
  final int currentFatigue;

  MonsterData({
    required this.name,
    required this.health,
    required this.fatigue,
    required this.currentHealth,
    required this.currentFatigue,
  });

  factory MonsterData.fromJson(Map<String, dynamic> json) {
    return MonsterData(
      name: json['name'],
      health: json['health'],
      fatigue: json['fatigue'],
      currentHealth: json['current_health'],
      currentFatigue: json['current_fatigue'],
    );
  }
}

class MonsterDetailedData {
  final String monsterName;
  final int monsterStrength;
  final int monsterDexterity;
  final int monsterIntelligence;
  final int monsterCurrentStrength;
  final int monsterCurrentDexterity;
  final int monsterCurrentIntelligence;
  final int monsterHealth;
  final int monsterFatigue;
  final int monsterCurrentHealth;
  final int monsterCurrentFatigue;
  final List<ObjectDetailedData>? monsterEquippedObjects;

  MonsterDetailedData({
    required this.monsterName,
    required this.monsterStrength,
    required this.monsterDexterity,
    required this.monsterIntelligence,
    required this.monsterCurrentStrength,
    required this.monsterCurrentDexterity,
    required this.monsterCurrentIntelligence,
    required this.monsterHealth,
    required this.monsterFatigue,
    required this.monsterCurrentHealth,
    required this.monsterCurrentFatigue,
    required this.monsterEquippedObjects,
  });

  factory MonsterDetailedData.fromJson(Map<String, dynamic> json) {
    List<dynamic>? equippedObjects = json['equipped_objects'];
    List<ObjectDetailedData>? equippedObjectData;
    if (equippedObjects != null) {
      equippedObjectData =
          equippedObjects.map((e) => ObjectDetailedData.fromJson(e)).toList();
    }

    return MonsterDetailedData(
      monsterName: json['name'],
      monsterStrength: json['strength'],
      monsterDexterity: json['dexterity'],
      monsterIntelligence: json['intelligence'],
      monsterCurrentStrength: json['current_strength'],
      monsterCurrentDexterity: json['current_dexterity'],
      monsterCurrentIntelligence: json['current_intelligence'],
      monsterHealth: json['health'],
      monsterFatigue: json['fatigue'],
      monsterCurrentHealth: json['current_health'],
      monsterCurrentFatigue: json['current_fatigue'],
      monsterEquippedObjects: equippedObjectData,
    );
  }

  Map<String, dynamic> toJson() => {
        'name': monsterName,
        'strength': monsterStrength,
        'dexterity': monsterDexterity,
        'intelligence': monsterIntelligence,
        'current_strength': monsterCurrentStrength,
        'current_dexterity': monsterCurrentDexterity,
        'current_intelligence': monsterCurrentIntelligence,
        'health': monsterHealth,
        'fatigue': monsterFatigue,
        'current_health': monsterCurrentHealth,
        'current_fatigue': monsterCurrentFatigue,
      };
}
