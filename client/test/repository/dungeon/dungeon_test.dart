import 'package:flutter_test/flutter_test.dart';

// Application
import 'package:go_mud_client/repository/dungeon/dungeon_repository.dart';

// Local Test Utilities
import '../../utility.dart';

void main() {
  test('DungeonRepository should ', () async {
    final repository = DungeonRepository(config: getConfig(), api: getAPI());
    expect(repository, isNotNull);

    final dungeons = await repository.getMany();
    expect(dungeons, isNotEmpty);
  });
}
