// Package
part of 'dungeon_repository.dart';

class DungeonRecord extends Equatable {
  final String id;
  final String name;
  final String description;

  const DungeonRecord({
    required this.id,
    required this.name,
    required this.description,
  });

  factory DungeonRecord.fromJson(Map<String, dynamic> json) {
    if (json.isEmpty) {
      throw RecordEmptyException('JSON data is empty');
    }
    return DungeonRecord(
      id: json['id'],
      name: json['name'],
      description: json['description'],
    );
  }

  Map<String, dynamic> toJson() => {
        'id': id,
        'name': name,
        'description': description,
      };

  @override
  List<Object> get props => [id, name, description];
}
