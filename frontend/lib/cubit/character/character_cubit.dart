import 'dart:async';
import 'package:bloc/bloc.dart';
import 'package:meta/meta.dart';
import 'package:equatable/equatable.dart';

// Application
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/repository/repository.dart';
import 'package:go_mud_client/logger.dart';

part 'character_state.dart';

const int maxAttributes = 36;
const int maxCharacters = 3;

class CharacterCubit extends Cubit<CharacterState> {
  final Map<String, String> config;
  final RepositoryCollection repositories;
  final DungeonActionCubit dungeonActionCubit;

  List<CharacterRecord>? characterRecords;
  CharacterRecord? characterRecord;

  // streamSubscription is listening to events from the dungeon action
  // cubit, specifically events that might require this cubit to refresh
  // and emit an updated state.
  StreamSubscription? streamSubscription;

  CharacterCubit({
    required this.config,
    required this.repositories,
    required this.dungeonActionCubit,
  }) : super(const CharacterStateInitial()) {
    streamSubscription?.cancel();
    streamSubscription = dungeonActionCubit.stream.listen((event) {
      final log = getLogger('CharacterCubit', 'dungeonActionCubit(listener)');
      log.warning('Dungeon action cubit emitted event $event');

      // TODO: 12-implement-death: Refresh character record when character
      // enters the dungeon, exits the dungeon, or is killed in the dungeon.
    });
  }

  void clearCharacter() {
    characterRecord = null;
    emit(const CharacterStateInitial());
  }

  bool canCreateCharacter() {
    if (characterRecords == null) {
      return true;
    }
    if (characterRecords != null && characterRecords!.length <= maxCharacters) {
      return true;
    }
    return false;
  }

  Future<void> initCreateCharacter() async {
    if (canCreateCharacter()) {
      emit(const CharacterStateCreate());
    }
  }

  Future<void> createCharacter(CreateCharacterRecord characterRecord) async {
    final log = getLogger('CharacterCubit', 'createCharacter');
    log.fine('Creating character $characterRecord');

    emit(const CharacterStateCreate());

    if (characterRecord.characterStrength +
            characterRecord.characterDexterity +
            characterRecord.characterIntelligence >
        maxAttributes) {
      String message = 'New character attributes exceeds maximum allowed';
      log.warning(message);
      emit(CharacterStateCreateError(
          characterRecord: characterRecord, message: message));
      return;
    }

    CharacterRecord? createdCharacterRecord;

    try {
      createdCharacterRecord =
          await repositories.characterRepository.createOne(characterRecord);
    } on DuplicateCharacterNameException {
      log.warning('Throwing character create error');
      emit(CharacterStateCreateError(
        characterRecord: characterRecord,
        message:
            'Character name ${characterRecord.characterName} has been taken.',
      ));
      return;
    } on RepositoryException catch (err) {
      log.warning('Throwing character create error ${err.message}');
      emit(CharacterStateCreateError(
          characterRecord: characterRecord, message: err.message));
      return;
    }

    if (createdCharacterRecord != null) {
      log.fine('Created character $createdCharacterRecord');
      this.characterRecord = createdCharacterRecord;
      characterRecords ??= [];
      characterRecords?.add(createdCharacterRecord);
      emit(CharacterStateSelected(characterRecord: createdCharacterRecord));
    }
  }

  Future<void> loadCharacters() async {
    final log = getLogger('CharacterCubit', 'loadCharacters');
    log.fine('Loading characters...');
    emit(const CharacterStateLoading());

    List<CharacterRecord>? characterRecords;

    try {
      characterRecords = await repositories.characterRepository.getMany();
    } catch (err) {
      emit(const CharacterStateLoadError());
      return;
    }

    emit(CharacterStateLoaded(characterRecords: characterRecords));
  }

  Future<void> loadCharacter(String characterID) async {
    final log = getLogger('CharacterCubit', 'loadCharacter');
    log.fine('Loading character ID $characterID');

    emit(const CharacterStateLoading());

    CharacterRecord? loadedCharacterRecord =
        await repositories.characterRepository.getOne(characterID);

    log.fine('Loaded character $loadedCharacterRecord');

    if (loadedCharacterRecord != null) {
      emit(CharacterStateSelected(characterRecord: loadedCharacterRecord));
    }
  }

  Future<void> refreshCharacter(String characterID) async {
    final log = getLogger('CharacterCubit', 'refreshCharacter');
    log.fine('Refreshing character ID $characterID');

    emit(const CharacterStateLoading());

    CharacterRecord? characterRecord =
        await repositories.characterRepository.getOne(characterID);

    if (characterRecord != null) {
      return selectCharacter(characterRecord);
    }
  }

  Future<void> selectCharacter(CharacterRecord characterRecord) async {
    this.characterRecord = characterRecord;

    emit(
      CharacterStateLoaded(
        characterRecords: characterRecords,
        currentCharacterRecord: characterRecord,
      ),
    );
  }
}
