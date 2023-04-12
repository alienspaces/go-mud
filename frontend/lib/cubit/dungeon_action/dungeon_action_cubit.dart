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

  DungeonActionCubit({required this.config, required this.repositories})
      : super(const DungeonActionStateInitial());

  Future<void> createAction(
      String dungeonID, String characterID, String command) async {
    final log = getLogger('DungeonActionCubit', 'createAction');
    log.fine('Creating dungeon action command >$command<');

    // TODO: This state is not used by any widgets to do anything
    // meaningful at the moment..
    emit(DungeonActionStateCreating(
      sentence: command,
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
      log.fine(
          'Previous record turn ${previousRec.actionTurnNumber} serial ${previousRec.actionSerialNumber}');
      log.fine('-- location ${previousRec.actionLocation.locationName}');
      if (previousRec.actionTargetLocation != null) {
        log.fine(
            '-- targetLocation ${previousRec.actionTargetLocation?.locationName}');
      }
      if (previousRec.actionTargetCharacter != null) {
        log.fine(
            '-- targetCharacter ${previousRec.actionTargetCharacter?.toJson()}');
      }
      if (previousRec.actionTargetMonster != null) {
        log.fine(
            '-- targetMonster ${previousRec.actionTargetMonster?.toJson()}');
      }
      if (previousRec.actionTargetObject != null) {
        log.fine('-- targetObject ${previousRec.actionTargetObject?.toJson()}');
      }
    }

    dungeonActionRecords.add(currentRec);

    // TODO: 9-implement-monster-actions: The following should possibly be
    // current = playerAction and previous = otherPlayerAndMonsterActions
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
    dungeonActionRecords = [];
  }

  /// Returns true when there are more actions in history to play
  bool playAction() {
    final log = getLogger('DungeonActionCubit', 'playAction');

    // TODO: Is this totally necessary? On the first action there is no
    // previous action, can we skip setting previous here and assign current?
    // What breaks if we change the following check to < 1 ?
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

    // TODO: 9-implement-monster-actions: We should probably split this out
    // into two different states, one for the player character and another
    // for actions which are other character and monster actions. The previous
    // action record is necessary for rendering the scrolling animation when
    // moving between rooms but shouldn't be necessary for other player or
    // monster action animations?
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
