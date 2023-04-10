import 'package:flutter_test/flutter_test.dart';

// Application
import 'package:go_mud_client/repository/dungeon_action/dungeon_action_repository.dart';

// Local Test Utilities
import '../../utility.dart';

void main() {
  test('DungeonActionRepository should', () async {
    final repository =
        DungeonActionRepository(config: getConfig(), api: getAPI());
    expect(repository, isNotNull,
        reason: 'DungeonActionRepository is not null');

    // Look

    final dungeonActionRecord =
        await repository.create(testDungeonID, testCharacterID, 'look');
    expect(
      dungeonActionRecord,
      isNotNull,
      reason:
          'DungeonActionRepository create "look" command returns a dungeon action record',
    );
    expect(
      dungeonActionRecord!.actionCommand,
      isNotNull,
      reason: 'DungeonActionRecord.command is not null',
    );
    expect(
      dungeonActionRecord.actionLocation,
      isNotNull,
      reason: 'DungeonActionRecord.location is not null',
    );
    expect(
      dungeonActionRecord.actionCharacter ?? dungeonActionRecord.actionMonster,
      isNotNull,
      reason:
          'DungeonActionRecord.character or DungeonActionRecord.monster is not null',
    );
  });
}
