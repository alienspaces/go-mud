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
      throw RecordEmptyException('DungeonRecord');
    }
    return DungeonRecord(
      dungeonID: json['id'],
      dungeonName: json['name'],
      dungeonDescription: json['description'],
    );
  }

  Map<String, dynamic> toJson() => {
        'id': dungeonID,
        'name': dungeonName,
        'description': dungeonDescription,
      };

  @override
  List<Object> get props => [
        dungeonID,
        dungeonName,
        dungeonDescription,
      ];
}
