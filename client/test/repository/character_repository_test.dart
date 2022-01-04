import 'package:flutter_test/flutter_test.dart';

// Application
import 'package:go_mud_client/repository/character/character_repository.dart';

// Local Test Utilities
import '../utility.dart';

void main() {
  test('CharacterRepository', () async {
    final repository = CharacterRepository(config: getConfig(), api: getAPI());
    expect(repository, isNotNull, reason: 'CharacterRepository is not null');

    final characters = await repository.getMany(testDungeonID);
    expect(
      characters,
      isNotEmpty,
      reason: 'CharacterRepository getMany response is not empty',
    );

    final character = await repository.getOne(testDungeonID, testCharacterID);
    expect(
      character,
      isNotNull,
      reason: 'CharacterRepository getOne response is not null',
    );
    expect(
      character!.id,
      isNotNull,
      reason: 'CharacterRepository getOne character.id is not null',
    );
    expect(
      character.name,
      isNotNull,
      reason: 'CharacterRepository getOne character.name is not null',
    );
    expect(
      character.strength,
      isNotNull,
      reason: 'CharacterRepository getOne character.strength is not null',
    );
    expect(
      character.dexterity,
      isNotNull,
      reason: 'CharacterRepository getOne character.dexterity is not null',
    );
    expect(
      character.intelligence,
      isNotNull,
      reason: 'CharacterRepository getOne character.intelligence is not null',
    );
    expect(
      character.currentStrength,
      isNotNull,
      reason:
          'CharacterRepository getOne character.currentStrength is not null',
    );
    expect(
      character.currentDexterity,
      isNotNull,
      reason:
          'CharacterRepository getOne character.currentDexterity is not null',
    );
    expect(
      character.currentIntelligence,
      isNotNull,
      reason:
          'CharacterRepository getOne character.currentIntelligence is not null',
    );
    expect(
      character.health,
      isNotNull,
      reason: 'CharacterRepository getOne character.health is not null',
    );
    expect(
      character.fatigue,
      isNotNull,
      reason: 'CharacterRepository getOne character.fatigue is not null',
    );
  });
}
