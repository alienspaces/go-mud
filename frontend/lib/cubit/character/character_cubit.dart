import 'dart:async';
import 'package:bloc/bloc.dart';
import 'package:meta/meta.dart';
import 'package:equatable/equatable.dart';

// Application
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/repository/repository.dart';
import 'package:go_mud_client/logger.dart';

part 'character_state.dart';

class CharacterCubit extends Cubit<CharacterState> {
  final Map<String, String> config;
  final RepositoryCollection repositories;
  final DungeonActionCubit dungeonActionCubit;

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
    streamSubscription = dungeonActionCubit.stream.listen((state) {
      final log = getLogger('CharacterCubit', 'dungeonActionCubit(listener)');
      log.warning('Dungeon action cubit emitted state');
      if (state is DungeonActionStateError) {
        log.warning('Dungeon action cubit emitted error event');
        if (characterRecord != null) {
          log.warning('Clearing character record');
          clearCharacter();
        }
      }
    });
  }

  void clearCharacter() {
    final log = getLogger('CharacterCubit', 'clearCharacter');
    log.warning('Clearing character');
    characterRecord = null;
    emit(const CharacterStateInitial());
  }

  Future<void> refresh() async {
    final log = getLogger('CharacterCubit', 'refresh');

    // Character not selected
    if (this.characterRecord == null) {
      // Exception here..
      return;
    }

    log.info('Refreshing character ID ${this.characterRecord!.characterID}');

    CharacterRecord? characterRecord = await repositories.characterRepository
        .getOne(this.characterRecord!.characterID);

    this.characterRecord = characterRecord;

    log.info('Refreshed character $characterRecord');
  }

  Future<void> select(CharacterRecord characterRecord) async {
    final log = getLogger('CharacterCubit', 'select');
    log.info('Selected character ID ${characterRecord.characterID}');
    log.info('Selected character Name ${characterRecord.characterName}');

    emit(CharacterStateSelecting(characterRecord: characterRecord));

    this.characterRecord = characterRecord;

    emit(CharacterStateSelected(characterRecord: characterRecord));
  }

  Future<void> enter(DungeonRecord dungeonRecord) async {
    final log = getLogger('CharacterCubit', 'enter');
    log.info('Entering dungeon ID ${dungeonRecord.dungeonID}');
    log.info('Entering dungeon Name ${dungeonRecord.dungeonName}');

    // Character not selected
    if (characterRecord == null) {
      // Exception here..
      return;
    }

    // Character already in this dungeon
    if (characterRecord != null &&
        characterRecord!.dungeonID == dungeonRecord.dungeonID) {
      // Exception here..
      return;
    }

    emit(CharacterStateEntering(characterRecord: characterRecord!));

    try {
      characterRecord = await repositories.dungeonCharacterRepository
          .enterDungeonCharacter(
              dungeonRecord.dungeonID, characterRecord!.characterID);
    } on RepositoryException catch (err) {
      log.warning('Throwing dungeon character enter error');
      emit(CharacterStateError(
          characterRecord: characterRecord!, message: err.message));
      return;
    }

    if (characterRecord != null) {
      log.fine('Entered dungeon with character $characterRecord');
      emit(CharacterStateEntered(
        characterRecord: characterRecord!,
      ));
    }
  }

  Future<void> exit() async {
    final log = getLogger('CharacterCubit', 'exit');

    // Character not selected
    if (characterRecord == null || characterRecord!.dungeonID == null) {
      // Exception here..
      return;
    }

    log.fine(
        'Exiting dungeon ID ${characterRecord!.dungeonID} character ID ${characterRecord!.characterID}');

    emit(CharacterStateExiting(characterRecord: characterRecord!));

    try {
      await repositories.dungeonCharacterRepository.exitDungeonCharacter(
        characterRecord!.dungeonID!,
        characterRecord!.characterID,
      );
    } on RepositoryException catch (err) {
      log.warning('Throwing dungeon character exit error');
      emit(CharacterStateError(
        characterRecord: characterRecord!,
        message: err.message,
      ));
      return;
    }

    log.fine('Exited dungeon character $characterRecord');
    emit(CharacterStateExited(characterRecord: characterRecord!));
  }
}
