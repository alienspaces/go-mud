import 'package:flutter_test/flutter_test.dart';

// Application
import 'package:go_mud_client/repository/dungeon/dungeon_repository.dart';

// Local Test Utilities
import '../../utility.dart';

void main() {
  test('DungeonRepository', () async {
    final repository = DungeonRepository(config: getConfig(), api: getAPI());
    expect(repository, isNotNull, reason: 'DungeonRepository is not null');

    final dungeons = await repository.getMany();
    expect(
      dungeons,
      isNotEmpty,
      reason: 'DungeonRepository getMany response is not empty',
    );

    final dungeon = await repository.getOne(testDungeonID);
    expect(
      dungeon,
      isNotNull,
      reason: 'DungeonRepository getOne response is not null',
    );
    expect(
      dungeon!.dungeonID,
      isNotNull,
      reason: 'DungeonRepository getOne dungeon.id is not null',
    );
    expect(
      dungeon.dungeonName,
      isNotNull,
      reason: 'DungeonRepository getOne dungeon.name is not null',
    );
    expect(
      dungeon.dungeonDescription,
      isNotNull,
      reason: 'DungeonRepository getOne dungeon.description is not null',
    );
  });
}
