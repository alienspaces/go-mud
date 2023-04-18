import 'package:flutter_test/flutter_test.dart';

// Application
import 'package:go_mud_client/repository/dungeon_action/dungeon_action_repository.dart';

// Local Test Utilities
import '../../utility.dart';

void main() {
  test('DungeonActionRepository', () async {
    final repository =
        DungeonActionRepository(config: getConfig(), api: getAPI());
    expect(repository, isNotNull,
        reason: 'DungeonActionRepository is not null');

    // Look
    final actionRecs =
        await repository.create(testDungeonID, testCharacterID, 'look');
    expect(
      actionRecs,
      isNotNull,
      reason:
          'var DungeonActionRepository create "look" command returns a dungeon action record',
    );

    for (var actionRec in actionRecs!) {
      expect(
        actionRec.actionCommand,
        isNotNull,
        reason: 'DungeonActionRecord.command is not null',
      );
      expect(
        actionRec.actionLocation,
        isNotNull,
        reason: 'DungeonActionRecord.location is not null',
      );
      expect(
        actionRec.actionCharacter ?? actionRec.actionMonster,
        isNotNull,
        reason:
            'DungeonActionRecord.character or DungeonActionRecord.monster is not null',
      );
    }
  });
}
