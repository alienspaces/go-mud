import 'package:flutter_test/flutter_test.dart';

// Application
import 'package:go_mud_client/repository/character/character_repository.dart';

// Local Test Utilities
import '../../utility.dart';

void main() {
  test('CharacterRepository', () async {
    final repository = CharacterRepository(config: getConfig(), api: getAPI());
    expect(repository, isNotNull, reason: 'CharacterRepository is not null');

    final characters = await repository.getMany();
    expect(
      characters,
      isNotEmpty,
      reason: 'CharacterRepository getMany response is not empty',
    );

    final character = await repository.getOne(testCharacterID);
    expect(
      character,
      isNotNull,
      reason: 'CharacterRepository getOne response is not null',
    );
    expect(
      character!.characterID,
      isNotNull,
      reason: 'CharacterRepository getOne character.id is not null',
    );
    expect(
      character.characterName,
      isNotNull,
      reason: 'CharacterRepository getOne character.name is not null',
    );
    expect(
      character.characterStrength,
      isNotNull,
      reason: 'CharacterRepository getOne character.strength is not null',
    );
    expect(
      character.characterDexterity,
      isNotNull,
      reason: 'CharacterRepository getOne character.dexterity is not null',
    );
    expect(
      character.characterIntelligence,
      isNotNull,
      reason: 'CharacterRepository getOne character.intelligence is not null',
    );
    expect(
      character.characterHealth,
      isNotNull,
      reason: 'CharacterRepository getOne character.health is not null',
    );
    expect(
      character.characterFatigue,
      isNotNull,
      reason: 'CharacterRepository getOne character.fatigue is not null',
    );
  });
}
