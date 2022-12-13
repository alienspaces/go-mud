// Package
part of 'dungeon_repository.dart';

class DungeonRecord extends Equatable {
  final String dungeonID;
  final String dungeonName;
  final String dungeonDescription;

  const DungeonRecord({
    required this.dungeonID,
    required this.dungeonName,
    required this.dungeonDescription,
  });

  factory DungeonRecord.fromJson(Map<String, dynamic> json) {
    if (json.isEmpty) {
      throw RecordEmptyException('JSON data is empty');
    }
    return DungeonRecord(
      dungeonID: json['dungeon_id'],
      dungeonName: json['dungeon_name'],
      dungeonDescription: json['dungeon_description'],
    );
  }

  Map<String, dynamic> toJson() => {
        'dungeon_id': dungeonID,
        'dungeon_name': dungeonName,
        'dungeon_description': dungeonDescription,
      };

  @override
  List<Object> get props => [
        dungeonID,
        dungeonName,
        dungeonDescription,
      ];
}
