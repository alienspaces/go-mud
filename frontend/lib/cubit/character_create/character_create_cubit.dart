import 'dart:async';
import 'package:bloc/bloc.dart';
import 'package:meta/meta.dart';
import 'package:equatable/equatable.dart';

// Application
import 'package:go_mud_client/repository/repository.dart';
import 'package:go_mud_client/logger.dart';

part 'character_create_state.dart';

const int maxAttributes = 36;
const int maxCharacters = 3;

class CharacterCreateCubit extends Cubit<CharacterCreateState> {
  final Map<String, String> config;
  final RepositoryCollection repositories;

  CharacterRecord? characterRecord;

  CharacterCreateCubit({
    required this.config,
    required this.repositories,
  }) : super(const CharacterCreateStateInitial());

  /// clearCharacter clears the current selected character and sets state
  /// back to the initial state.
  void clearCharacter() {
    final log = getLogger('CharacterCreateCubit', 'clearCharacter');
    log.warning('Clearing character');
    characterRecord = null;
    emit(const CharacterCreateStateInitial());
  }

  Future<void> createCharacter(CreateCharacterRecord characterRecord) async {
    final log = getLogger('CharacterCreateCubit', 'createCharacter');
    log.warning('Creating character $characterRecord');

    emit(const CharacterCreateStateCreating());

    if (characterRecord.characterStrength +
            characterRecord.characterDexterity +
            characterRecord.characterIntelligence >
        maxAttributes) {
      String message = 'New character attributes exceeds maximum allowed';
      log.warning(message);
      emit(CharacterCreateStateError(
          characterRecord: characterRecord, message: message));
      return;
    }

    CharacterRecord? createdCharacterRecord;

    try {
      createdCharacterRecord =
          await repositories.characterRepository.createOne(characterRecord);
    } on DuplicateCharacterNameException {
      log.warning('Throwing character create error');
      emit(CharacterCreateStateError(
        characterRecord: characterRecord,
        message:
            'Character name ${characterRecord.characterName} has been taken.',
      ));
      return;
    } on RepositoryException catch (err) {
      log.warning('Throwing character create error ${err.message}');
      emit(CharacterCreateStateError(
          characterRecord: characterRecord, message: err.message));
      return;
    }

    if (createdCharacterRecord != null) {
      log.fine('Created character $createdCharacterRecord');
      this.characterRecord = createdCharacterRecord;
      emit(
          CharacterCreateStateCreated(characterRecord: createdCharacterRecord));
    }
  }
}
