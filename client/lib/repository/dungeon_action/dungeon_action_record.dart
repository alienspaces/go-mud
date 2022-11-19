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
    Map<String, dynamic>? location = json['action_location'];
    if (location == null) {
      throw const FormatException('Missing "action_location" from JSON');
    }

    var locationData = LocationData.fromJson(location);

    // Character or Monster
    Map<String, dynamic>? character = json['action_character'];
    Map<String, dynamic>? monster = json['action_monster'];
    if (character == null && monster == null) {
      throw const FormatException(
          'Missing "action_character" or "action_monster" from JSON');
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
    Map<String, dynamic>? equippedObject = json['action_equipped_object'];
    ObjectDetailedData? equippedObjectData;
    if (equippedObject != null) {
      equippedObjectData = ObjectDetailedData.fromJson(equippedObject);
    }

    // Stashed object
    Map<String, dynamic>? stashedObject = json['action_stashed_object'];
    ObjectDetailedData? stashedObjectData;
    if (stashedObject != null) {
      stashedObjectData = ObjectDetailedData.fromJson(stashedObject);
    }

    // Dropped object
    Map<String, dynamic>? droppedObject = json['action_dropped_object'];
    ObjectDetailedData? droppedObjectData;
    if (droppedObject != null) {
      droppedObjectData = ObjectDetailedData.fromJson(droppedObject);
    }

    // Target object
    Map<String, dynamic>? targetObject = json['action_target_object'];
    ObjectDetailedData? targetObjectData;
    if (targetObject != null) {
      targetObjectData = ObjectDetailedData.fromJson(targetObject);
    }

    // Target character
    Map<String, dynamic>? targetCharacter = json['action_target_character'];
    CharacterDetailedData? targetCharacterData;
    if (targetCharacter != null) {
      targetCharacterData = CharacterDetailedData.fromJson(targetCharacter);
    }

    // Target monster
    Map<String, dynamic>? targetMonster = json['action_target_monster'];
    MonsterDetailedData? targetMonsterData;
    if (targetMonster != null) {
      targetMonsterData = MonsterDetailedData.fromJson(targetMonster);
    }

    // Target location
    LocationData? targetLocationData;

    Map<String, dynamic>? targetLocation = json['action_target_location'];
    if (targetLocation != null) {
      targetLocationData = LocationData.fromJson(targetLocation);
    }

    return DungeonActionRecord(
      actionID: json['action_id'],
      actionCommand: json['action_command'],
      actionNarrative: json['action_narrative'],
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
    List<dynamic> directions = json['location_directions'];

    // Location objects
    List<dynamic>? locationObjects = json['location_objects'];
    List<ObjectData>? locationObjectData;
    if (locationObjects != null) {
      locationObjectData =
          locationObjects.map((e) => ObjectData.fromJson(e)).toList();
    }

    // Location characters
    List<dynamic>? locationCharacters = json['location_characters'];
    List<CharacterData>? locationCharacterData;
    if (locationCharacters != null) {
      locationCharacterData =
          locationCharacters.map((e) => CharacterData.fromJson(e)).toList();
    }

    // Location monsters
    List<dynamic>? locationMonsters = json['location_monsters'];
    List<MonsterData>? locationMonsterData;
    if (locationMonsters != null) {
      locationMonsterData =
          locationMonsters.map((e) => MonsterData.fromJson(e)).toList();
    }

    return LocationData(
      locationName: json['location_name'],
      locationDescription: json['location_description'],
      locationDirection: json['location_direction'],
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
      objectName: json['object_name'],
      objectDescription: json['object_description'],
      objectIsStashed: json['object_is_stashed'],
      objectIsEquipped: json['object_is_equipped'],
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
    List<dynamic>? equippedObjects = json['character_equipped_objects'];
    List<ObjectDetailedData>? equippedObjectData;
    if (equippedObjects != null) {
      equippedObjectData =
          equippedObjects.map((e) => ObjectDetailedData.fromJson(e)).toList();
    }

    List<dynamic>? stashedObjects = json['character_stashed_objects'];
    List<ObjectDetailedData>? stashedObjectData;
    if (stashedObjects != null) {
      stashedObjectData =
          stashedObjects.map((e) => ObjectDetailedData.fromJson(e)).toList();
    }

    return CharacterDetailedData(
      characterName: json['character_name'],
      characterStrength: json['character_strength'],
      characterDexterity: json['character_dexterity'],
      characterIntelligence: json['character_intelligence'],
      characterCurrentStrength: json['character_current_strength'],
      characterCurrentDexterity: json['character_current_dexterity'],
      characterCurrentIntelligence: json['character_current_intelligence'],
      characterHealth: json['character_health'],
      characterFatigue: json['character_fatigue'],
      characterCurrentHealth: json['character_current_health'],
      characterCurrentFatigue: json['character_current_fatigue'],
      characterStashedObjects: stashedObjectData,
      characterEquippedObjects: equippedObjectData,
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
    List<dynamic>? equippedObjects = json['monster_equipped_objects'];
    List<ObjectDetailedData>? equippedObjectData;
    if (equippedObjects != null) {
      equippedObjectData =
          equippedObjects.map((e) => ObjectDetailedData.fromJson(e)).toList();
    }

    return MonsterDetailedData(
      monsterName: json['monster_name'],
      monsterStrength: json['monster_strength'],
      monsterDexterity: json['monster_dexterity'],
      monsterIntelligence: json['monster_intelligence'],
      monsterCurrentStrength: json['monster_current_strength'],
      monsterCurrentDexterity: json['monster_current_dexterity'],
      monsterCurrentIntelligence: json['monster_current_intelligence'],
      monsterHealth: json['monster_health'],
      monsterFatigue: json['monster_fatigue'],
      monsterCurrentHealth: json['monster_current_health'],
      monsterCurrentFatigue: json['monster_current_fatigue'],
      monsterEquippedObjects: equippedObjectData,
    );
  }
}
