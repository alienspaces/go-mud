import 'package:go_mud_client/api/api.dart';

// Application
import 'package:go_mud_client/repository/character/character_repository.dart';
import 'package:go_mud_client/repository/dungeon/dungeon_repository.dart';
import 'package:go_mud_client/repository/dungeon_action/dungeon_action_repository.dart';

class RepositoryCollection {
  final Map<String, String> config;
  final API api;
  late final CharacterRepository characterRepository;
  late final DungeonRepository dungeonRepository;
  late final DungeonActionRepository dungeonActionRepository;

  RepositoryCollection({required this.config, required this.api}) {
    characterRepository = CharacterRepository(config: config, api: api);
    dungeonRepository = DungeonRepository(config: config, api: api);
    dungeonActionRepository = DungeonActionRepository(config: config, api: api);
  }
}
