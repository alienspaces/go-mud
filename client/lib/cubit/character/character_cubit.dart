import 'package:bloc/bloc.dart';
import 'package:meta/meta.dart';
import 'package:equatable/equatable.dart';

// Application
import 'package:go_mud_client/repository/repository.dart';
import 'package:go_mud_client/logger.dart';

part 'character_state.dart';

const int maxAttributes = 36;

class CharacterCubit extends Cubit<CharacterState> {
  final Map<String, String> config;
  final RepositoryCollection repositories;

  CharacterRecord? characterRecord;

  CharacterCubit({required this.config, required this.repositories})
      : super(const CharacterStateInitial());

  void clearCharacter() {
    characterRecord = null;
    emit(const CharacterStateInitial());
  }

  Future<void> createCharacter(
      String dungeonID, CreateCharacterRecord characterRecord) async {
    final log = getLogger('CharacterCubit');
    log.fine('Creating character $characterRecord');

    emit(const CharacterStateCreating());

    if (characterRecord.strength +
            characterRecord.dexterity +
            characterRecord.intelligence >
        maxAttributes) {
      String message = 'New character attributes exceeds maximum allowed';
      log.warning(message);
      emit(CharacterStateCreateError(
          characterRecord: characterRecord, message: message));
      return;
    }

    CharacterRecord? createdCharacterRecord;

    try {
      createdCharacterRecord = await repositories.characterRepository
          .createOne(dungeonID, characterRecord);
    } on RepositoryException catch (err) {
      log.warning('Throwing character create error');
      emit(CharacterStateCreateError(
          characterRecord: characterRecord, message: err.message));
      return;
    }

    if (createdCharacterRecord != null) {
      log.fine('Created character $createdCharacterRecord');
      this.characterRecord = createdCharacterRecord;
      emit(CharacterStateSelected(characterRecord: createdCharacterRecord));
    }
  }

  Future<void> loadCharacter(String dungeonID, String characterID) async {
    final log = getLogger('CharacterCubit');
    log.fine('Creating character ID $characterID');

    emit(const CharacterStateLoading());

    CharacterRecord? loadedCharacterRecord =
        await repositories.characterRepository.getOne(dungeonID, characterID);

    log.fine('Created character $loadedCharacterRecord');

    if (loadedCharacterRecord != null) {
      emit(CharacterStateSelected(characterRecord: loadedCharacterRecord));
    }
  }
}
