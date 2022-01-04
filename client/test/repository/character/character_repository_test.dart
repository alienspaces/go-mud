import 'package:flutter_test/flutter_test.dart';

// Application
import 'package:go_mud_client/repository/character/character_repository.dart';

// Local Test Utilities
import '../../utility.dart';

void main() {
  test('CharacterRepository should ', () async {
    final repository = CharacterRepository(config: getConfig(), api: getAPI());
    expect(repository, isNotNull, reason: "CharacterRepository is not null");

    final characters = await repository.getMany(testDungeonID);
    expect(
      characters,
      isNotEmpty,
      reason: "CharacterRepository getMany response is not empty",
    );
  });
}
