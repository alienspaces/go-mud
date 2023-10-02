import 'dart:async';
import 'package:bloc/bloc.dart';
import 'package:meta/meta.dart';
import 'package:equatable/equatable.dart';

// Application
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/repository/repository.dart';
import 'package:go_mud_client/logger.dart';

part 'dungeon_character_state.dart';

class DungeonCharacterCubit extends Cubit<DungeonCharacterState> {
  final Map<String, String> config;
  final RepositoryCollection repositories;
  final DungeonActionCubit dungeonActionCubit;

  // dungeonCharacterRecord will be populated when there is a character
  // being played and they have entered into a dungeon instance.
  DungeonCharacterRecord? dungeonCharacterRecord;

  // streamSubscription is listening to events from the dungeon action
  // cubit, specifically events that might require this cubit to refresh
  // and emit an updated state.
  StreamSubscription? streamSubscription;

  DungeonCharacterCubit({
    required this.config,
    required this.repositories,
    required this.dungeonActionCubit,
  }) : super(const DungeonCharacterStateInitial()) {
    streamSubscription?.cancel();
    streamSubscription = dungeonActionCubit.stream.listen((event) {
      final log = getLogger('CharacterCubit', 'dungeonActionCubit(listener)');

      // TODO: 12-implement-death: Refresh character record when character
      // enters the dungeon, exits the dungeon, or is killed in the dungeon.
      // Continue testing what this does when a character dies, we want to go
      // back to the character selection screen.
      if (state is DungeonActionStateError) {
        log.warning('Dungeon action cubit emitted error event');
        if (dungeonCharacterRecord != null) {
          log.warning('Clearing dungeon character record');
          clearDungeonCharacter();
        }
      }
    });
  }

  void clearDungeonCharacter() {
    final log = getLogger('CharacterCubit', 'clearCharacter');
    log.warning('Clearing character');
    dungeonCharacterRecord = null;
    emit(const DungeonCharacterStateInitial());
  }

  /// getDungeonCharacterRecordForCharacter returns a DungeonCharacterRecord if
  /// the character is currently inside a dungeon
  Future<DungeonCharacterRecord?> getDungeonCharacterRecord(
      String dungeonID, String characterID) async {
    final log = getLogger('DungeonCharacterCubit', 'getDungeonCharacterRecord');

    emit(DungeonCharacterStateLoading(characterID: characterID));

    DungeonCharacterRecord? rec;
    try {
      rec = await repositories.dungeonCharacterRepository
          .getDungeonCharacter(dungeonID, characterID);
    } on RepositoryException catch (err) {
      log.warning('Throwing dungeon character load error');
      emit(DungeonCharacterStateLoadError(
        characterID: characterID,
        message: err.message,
      ));
      return Future<DungeonCharacterRecord?>.value(null);
    }

    return Future<DungeonCharacterRecord?>.value(rec);
  }

  Future<void> enterDungeonCharacter(
      String dungeonID, String characterID) async {
    final log = getLogger('DungeonCharacterCubit', 'enterDungeonCharacter');
    log.fine('Entering dungeon ID $dungeonID character ID $characterID');

    emit(const DungeonCharacterStateCreate());

    dungeonCharacterRecord =
        await getDungeonCharacterRecord(dungeonID, characterID);

    if (dungeonCharacterRecord != null &&
        dungeonCharacterRecord?.dungeonID == dungeonID &&
        dungeonCharacterRecord?.characterID == characterID) {
      log.fine(
          'Dungeon with character $dungeonCharacterRecord is already in this dungeon, resuming..');
      emit(DungeonCharacterStateCreated(
        dungeonCharacterRecord: dungeonCharacterRecord!,
      ));
      return;
    }

    // Character already inside some other dungeon
    if (dungeonCharacterRecord != null) {
      log.fine(
          'Dungeon with character $dungeonCharacterRecord is already in a dungeon');
      emit(DungeonCharacterStateCreateError(
          dungeonID: dungeonID,
          characterID: characterID,
          message:
              'Dungeon with character $dungeonCharacterRecord is already in a dungeon'));
      return;
    }

    try {
      dungeonCharacterRecord = await repositories.dungeonCharacterRepository
          .enterDungeonCharacter(dungeonID, characterID);
    } on RepositoryException catch (err) {
      log.warning('Throwing dungeon character enter error');
      emit(DungeonCharacterStateCreateError(
          dungeonID: dungeonID,
          characterID: characterID,
          message: err.message));
      return;
    }

    if (dungeonCharacterRecord != null) {
      log.fine('Entered dungeon with character $dungeonCharacterRecord');
      emit(DungeonCharacterStateCreated(
        dungeonCharacterRecord: dungeonCharacterRecord!,
      ));
    }
  }

  Future<void> exitDungeonCharacter(
      String dungeonID, String characterID) async {
    final log = getLogger('DungeonCharacterCubit', 'exitDungeonCharacter');
    log.fine('Exiting dungeon ID $dungeonID character ID $characterID');

    emit(const DungeonCharacterStateDelete());

    DungeonCharacterRecord? dungeonCharacterRecord =
        await getDungeonCharacterRecord(dungeonID, characterID);

    if (dungeonCharacterRecord == null) {
      log.warning(
          'Did not find dungeon character record for dungeon ID $dungeonID character ID $characterID');

      emit(DungeonCharacterStateDeleteError(
          dungeonID: dungeonID,
          characterID: characterID,
          message: 'Exiting dungeon character record does not exist'));
      return;
    }

    try {
      await repositories.dungeonCharacterRepository
          .exitDungeonCharacter(dungeonID, characterID);
    } on RepositoryException catch (err) {
      log.warning('Throwing dungeon character exit error');
      emit(DungeonCharacterStateDeleteError(
          dungeonID: dungeonID,
          characterID: characterID,
          message: err.message));
      return;
    }

    log.fine('Exited dungeon character $dungeonCharacterRecord');
    emit(DungeonCharacterStateDeleted(
        dungeonCharacterRecord: dungeonCharacterRecord));
    dungeonCharacterRecord = null;
  }
}
