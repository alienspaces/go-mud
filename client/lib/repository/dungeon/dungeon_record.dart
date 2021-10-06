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

  DungeonRecord.fromJson(Map<String, dynamic> json)
      : id = json['id'],
        name = json['name'],
        description = json['description'];

  Map<String, dynamic> toJson() => {
        'id': id,
        'name': name,
        'description': description,
      };

  @override
  List<Object> get props => [id, name, description];
}
