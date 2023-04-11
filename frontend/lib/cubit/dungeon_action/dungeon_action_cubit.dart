import 'package:bloc/bloc.dart';
import 'package:meta/meta.dart';
import 'package:equatable/equatable.dart';

// Application
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/repository/repository.dart';

part 'dungeon_action_state.dart';

class DungeonActionCubit extends Cubit<DungeonActionState> {
  final Map<String, String> config;
  final RepositoryCollection repositories;

  List<DungeonActionRecord> dungeonActionRecords = [];
  DungeonActionRecord? dungeonActionRecord;

  DungeonActionCubit({required this.config, required this.repositories})
      : super(const DungeonActionStateInitial());

  Future<void> createAction(
      String dungeonID, String characterID, String command) async {
    final log = getLogger('DungeonActionCubit', 'createAction');
    log.fine('Creating dungeon action command >$command<');

    emit(DungeonActionStateCreating(
      sentence: command,
      current: dungeonActionRecord,
    ));

    List<DungeonActionRecord>? createdRecs;
    try {
      createdRecs = await repositories.dungeonActionRepository
          .create(dungeonID, characterID, command);
    } on ActionTooEarlyException {
      emit(const DungeonActionStateError(
        message: 'Command too early, slow down...',
      ));
      return;
    } on RepositoryException catch (err) {
      log.warning('Throwing action state error ${err.message}');
      emit(DungeonActionStateError(
        message: err.message,
      ));
      return;
    }

    if (createdRecs == null || createdRecs.isEmpty) {
      log.warning('No action records returned');
      return;
    }

    DungeonActionRecord currentRec = createdRecs.removeLast();
    List<DungeonActionRecord>? previousRecs = createdRecs;

    for (var previousRec in previousRecs) {
      log.info(
          'Previous record turn ${previousRec.actionTurnNumber} serial ${previousRec.actionSerialNumber}');
      log.info('-- location ${previousRec.actionLocation.locationName}');
      if (previousRec.actionTargetLocation != null) {
        log.info(
            '-- targetLocation ${previousRec.actionTargetLocation?.locationName}');
      }
      if (previousRec.actionTargetCharacter != null) {
        log.info(
            '-- targetCharacter ${previousRec.actionTargetCharacter?.toJson()}');
      }
      if (previousRec.actionTargetMonster != null) {
        log.info(
            '-- targetMonster ${previousRec.actionTargetMonster?.toJson()}');
      }
      if (previousRec.actionTargetObject != null) {
        log.info('-- targetObject ${previousRec.actionTargetObject?.toJson()}');
      }
    }

    emit(
      DungeonActionStateCreated(
        current: currentRec,
        previous: previousRecs,
        action: currentRec.actionCommand,
        direction: currentRec.actionTargetLocation?.locationDirection,
      ),
    );
  }

  /// Clear dungeon action history, typically called when returning to the dungeon page
  void clearActions() {
    dungeonActionRecord = null;
    dungeonActionRecords = [];
  }

  /// Returns true when there are more actions in history to play
  bool playAction() {
    final log = getLogger('DungeonActionCubit', 'playAction');

    if (dungeonActionRecords.length < 2) {
      log.info('Not enough dungeon action records, not playing action');
      return false;
    }

    DungeonActionRecord previous = dungeonActionRecords.removeAt(0);
    DungeonActionRecord current = dungeonActionRecords[0];
    String? direction;
    if (current.actionCommand == 'move' || current.actionCommand == 'look') {
      direction = current.actionTargetLocation?.locationDirection;
    }

    if (current.actionTargetLocation != null) {
      log.fine(
          'Play action command >${current.actionCommand}< direction >$direction<');
    }

    if (current.actionTargetCharacter != null) {
      log.fine(
          'Play action command >${current.actionCommand}< character >${current.actionTargetCharacter?.characterName}<');
    }

    if (current.actionTargetMonster != null) {
      log.fine(
          'Play action command >${current.actionCommand}< monster >${current.actionTargetMonster?.monsterName}<');
    }

    if (current.actionTargetObject != null) {
      log.fine(
          'Play action command >${current.actionCommand}< object >${current.actionTargetObject?.objectName}<');
    }

    emit(
      DungeonActionStatePlaying(
        previous: previous,
        current: current,
        action: current.actionCommand,
        direction: direction,
      ),
    );

    if (dungeonActionRecords.length <= 1) {
      return false;
    }

    return true;
  }
}
