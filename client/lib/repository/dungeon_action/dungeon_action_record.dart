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
  final DungeonAction action;
  final DungeonLocation location;
  final List<DungeonObject>? objects;
  final List<DungeonCharacter>? characters;
  final List<DungeonMonster>? monsters;

  const DungeonActionRecord({
    required this.action,
    required this.location,
    required this.objects,
    required this.characters,
    required this.monsters,
  });

  factory DungeonActionRecord.fromJson(Map<String, dynamic> json) {
    DungeonAction? dungeonAction;
    Map<String, dynamic>? action = json['action'];
    if (action == null) {
      throw const FormatException('Missing "action" from JSON');
    }
    dungeonAction = DungeonAction(
      id: action['id'],
      command: action['command'],
      commandResult: action['command_result'],
      equippedDungeonObjectName: action['equipped_dungeon_object_name'],
      stashedDungeonObjectName: action['stashed_dungeon_object_name'],
      targetDungeonObjectName: action['target_dungeon_object_name'],
      targetDungeonCharacterName: action['target_dungeon_character_name'],
      targetDungeonMonsterName: action['target_dungeon_monster_name'],
      targetDungeonLocationDirection: action['target_dungeon_location_direction'],
      targetDungeonLocationName: action['target_dungeon_location_name'],
    );

    DungeonLocation? dungeonLocation;
    Map<String, dynamic>? location = json['location'];
    if (location == null) {
      throw const FormatException('Missing "location" from JSON');
    }

    List<dynamic> directions = location['directions'];
    dungeonLocation = DungeonLocation(
      name: location['name'],
      description: location['description'],
      directions: directions.map((e) => e.toString()).toList(),
    );

    List<dynamic>? objects = json['objects'];
    List<DungeonObject>? dungeonObjects;
    if (objects != null) {
      dungeonObjects = objects.map((e) => DungeonObject.fromJson(e)).toList();
    }

    List<dynamic>? characters = json['characters'];
    List<DungeonCharacter>? dungeonCharacters;
    if (characters != null) {
      dungeonCharacters = characters.map((e) => DungeonCharacter.fromJson(e)).toList();
    }

    List<dynamic>? monsters = json['monsters'];
    List<DungeonMonster>? dungeonMonsters;
    if (monsters != null) {
      dungeonMonsters = monsters.map((e) => DungeonMonster.fromJson(e)).toList();
    }

    return DungeonActionRecord(
      action: dungeonAction,
      location: dungeonLocation,
      objects: dungeonObjects,
      characters: dungeonCharacters,
      monsters: dungeonMonsters,
    );
  }

  @override
  List<Object> get props => [
        action,
        location,
      ];
}

class DungeonAction {
  final String id;
  final String command;
  final String? commandResult;
  final String? equippedDungeonObjectName;
  final String? stashedDungeonObjectName;
  final String? targetDungeonObjectName;
  final String? targetDungeonCharacterName;
  final String? targetDungeonMonsterName;
  final String? targetDungeonLocationDirection;
  final String? targetDungeonLocationName;

  DungeonAction({
    required this.id,
    required this.command,
    this.commandResult,
    this.equippedDungeonObjectName,
    this.stashedDungeonObjectName,
    this.targetDungeonObjectName,
    this.targetDungeonCharacterName,
    this.targetDungeonMonsterName,
    this.targetDungeonLocationDirection,
    this.targetDungeonLocationName,
  });
}

class DungeonLocation {
  final String name;
  final String description;
  final List<String> directions;
  DungeonLocation({required this.name, required this.description, required this.directions});
}

class DungeonObject {
  final String name;
  DungeonObject({required this.name});

  factory DungeonObject.fromJson(Map<String, dynamic> json) {
    return DungeonObject(name: json['name']);
  }
}

class DungeonCharacter {
  final String name;
  DungeonCharacter({required this.name});

  factory DungeonCharacter.fromJson(Map<String, dynamic> json) {
    return DungeonCharacter(name: json['name']);
  }
}

class DungeonMonster {
  final String name;
  DungeonMonster({required this.name});

  factory DungeonMonster.fromJson(Map<String, dynamic> json) {
    return DungeonMonster(name: json['name']);
  }
}
