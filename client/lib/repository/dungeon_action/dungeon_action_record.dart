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

  const DungeonActionRecord({
    required this.action,
    required this.location,
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
      commandResult: action['commandResult'],
      equippedDungeonObjectName: action['equippedDungeonObjectName'],
      stashedDungeonObjectName: action['stashedDungeonObjectName'],
      targetDungeonObjectName: action['targetDungeonObjectName'],
      targetDungeonCharacterName: action['targetDungeonCharacterName'],
      targetDungeonMonsterName: action['targetDungeonMonsterName'],
      targetDungeonLocationDirection: action['targetDungeonLocationDirection'],
      targetDungeonLocationName: action['targetDungeonLocationName'],
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

    return DungeonActionRecord(
      action: dungeonAction,
      location: dungeonLocation,
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
