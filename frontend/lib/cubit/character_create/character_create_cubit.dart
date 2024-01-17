import 'dart:async';
import 'package:bloc/bloc.dart';
import 'package:go_mud_client/cubit/exception.dart';
import 'package:meta/meta.dart';
import 'package:equatable/equatable.dart';

// Application
import 'package:go_mud_client/repository/repository.dart';
import 'package:go_mud_client/logger.dart';

part 'character_create_state.dart';

const int maxAttributes = 36;
const int maxCharacters = 3;

class CreateCharacterResult {
  CubitException? exception;
  CharacterRecord? characterRecord;
}

class CharacterCreateCubit extends Cubit<CharacterCreateState> {
  final Map<String, String> config;
  final RepositoryCollection repositories;

  CharacterRecord? characterRecord;

  CharacterCreateCubit({
    required this.config,
    required this.repositories,
  }) : super(const CharacterCreateStateInitial());

  void clearCharacter() {
    final log = getLogger('CharacterCreateCubit', 'clearCharacter');
    log.warning('Clearing character');
    characterRecord = null;
    emit(const CharacterCreateStateInitial());
  }

  Future<CreateCharacterResult> createCharacter(
      CreateCharacterRecord characterRecord) async {
    final log = getLogger('CharacterCreateCubit', 'createCharacter');
    log.warning('Creating character $characterRecord');

    var result = CreateCharacterResult();

    emit(const CharacterCreateStateCreating());

    // Character attrbiutes exceed maximum
    if (characterRecord.characterStrength +
            characterRecord.characterDexterity +
            characterRecord.characterIntelligence >
        maxAttributes) {
      String message = 'New character attributes exceeds maximum allowed';
      log.warning(message);

      emit(CharacterCreateStateError(
        characterRecord: characterRecord,
        exception: CubitException(message),
      ));

      result.exception = CubitException(message);
      return result;
    }

    CharacterRecord? createdCharacterRecord;

    try {
      createdCharacterRecord =
          await repositories.characterRepository.createOne(characterRecord);
    }

    // Duplicate character error
    on DuplicateCharacterNameException {
      String message =
          'Character name ${characterRecord.characterName} has been taken.';
      log.warning('Throwing character create error');

      emit(CharacterCreateStateError(
        characterRecord: characterRecord,
        exception: CubitException(message),
      ));

      result.exception = CubitException(message);
      return result;
    }

    // Other error
    on RepositoryException catch (err) {
      String message = err.message;
      log.warning('Throwing character create error $message');

      emit(CharacterCreateStateError(
        characterRecord: characterRecord,
        exception: CubitException(message),
      ));

      result.exception = CubitException(message);
      return result;
    }

    // Character not created
    if (createdCharacterRecord == null) {
      result.exception = const CubitException('Character not created');
      return result;
    }

    // Character created
    log.fine('Created character $createdCharacterRecord');
    this.characterRecord = createdCharacterRecord;

    emit(
      CharacterCreateStateCreated(characterRecord: createdCharacterRecord),
    );

    result.characterRecord = createdCharacterRecord;
    return result;
  }
}
