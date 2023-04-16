import 'package:bloc/bloc.dart';
import 'package:meta/meta.dart';
import 'package:equatable/equatable.dart';
import 'package:logging/logging.dart';

// Application
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/repository/repository.dart';

part 'dungeon_action_state.dart';

class DungeonActionCubit extends Cubit<DungeonActionState> {
  final Map<String, String> config;
  final RepositoryCollection repositories;

  DungeonActionRecord? currentCharacterActionRec;
  DungeonActionRecord? previousCharacterActionRec;
  List<DungeonActionRecord> otherActionRecs = [];

  DungeonActionCubit({required this.config, required this.repositories})
      : super(const DungeonActionStateInitial());

  Future<void> createAction(
    String dungeonID,
    String characterID,
    String command,
  ) async {
    final log = getLogger('DungeonActionCubit', 'createAction');
    log.fine('Creating dungeon action command >$command<');

    // TODO: Not used by any widgets to do anything meaningful at the moment..
    emit(DungeonActionStateCreating(
      sentence: command,
    ));

    List<DungeonActionRecord>? responseActionRecs;
    try {
      responseActionRecs = await repositories.dungeonActionRepository
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

    if (responseActionRecs == null || responseActionRecs.isEmpty) {
      log.warning('No action records returned');
      return;
    }

    // The previous character action record is used for animating transitions
    // from the previous character action to the current character action
    previousCharacterActionRec = currentCharacterActionRec;

    // The first action record is always the character action
    currentCharacterActionRec = responseActionRecs.removeLast();

    // The remaining action records are always other character or monster actions
    List<DungeonActionRecord>? previousOtherActionRecs = responseActionRecs;

    while (previousOtherActionRecs.isNotEmpty) {
      // Add from the earliest other action
      var actionRec = previousOtherActionRecs.removeAt(0);
      logActionRec(log, actionRec);
      otherActionRecs.add(actionRec);
    }

    emit(
      DungeonActionStateCreated(
        action: currentCharacterActionRec!,
        previousAction: previousCharacterActionRec,
        actionCommand: currentCharacterActionRec!.actionCommand,
        direction:
            currentCharacterActionRec!.actionTargetLocation?.locationDirection,
      ),
    );
  }

  /// Clear dungeon action history, typically called when returning to the dungeon page
  void clearActions() {
    currentCharacterActionRec = null;
    previousCharacterActionRec = null;
    otherActionRecs = [];
  }

  /// Plays the current player action
  Future<void> playCharacterAction() async {
    final log = getLogger('DungeonActionCubit', 'playCharacterAction');

    if (currentCharacterActionRec == null) {
      log.warning(
          'Current character action is null, cannot play character action');
      return;
    }

    DungeonActionRecord actionRec = currentCharacterActionRec!;

    String? actionDirection;
    if (actionRec.actionCommand == 'move' ||
        actionRec.actionCommand == 'look') {
      actionDirection = actionRec.actionTargetLocation?.locationDirection;
    }

    await logActionRec(log, actionRec);

    emit(
      DungeonActionStatePlaying(
        currentActionRec: actionRec,
        previousActionRec: previousCharacterActionRec,
        actionCommand: actionRec.actionCommand,
        actionDirection: actionDirection,
      ),
    );

    return;
  }

  // Plays all existing other actions until none are left to play
  Future<void> playOtherActions() async {
    final log = getLogger('DungeonActionCubit', 'playOtherActions');

    if (otherActionRecs.isEmpty) {
      log.info('Other actions are empty, cannot play other action');
      return;
    }

    while (otherActionRecs.isNotEmpty) {
      // Play from the earliest other action
      DungeonActionRecord actionRec = otherActionRecs.removeAt(0);

      await logActionRec(log, actionRec);

      String? actionDirection;
      if (actionRec.actionCommand == 'move' ||
          actionRec.actionCommand == 'look') {
        actionDirection = actionRec.actionTargetLocation?.locationDirection;
      }

      emit(
        DungeonActionStatePlayingOther(
          actionRec: actionRec,
          actionCommand: actionRec.actionCommand,
          actionDirection: actionDirection,
        ),
      );
    }

    return;
  }
}

Future<void> logActionRec(Logger log, DungeonActionRecord actionRec) async {
  String? direction;

  if ((actionRec.actionCommand == 'move' ||
          actionRec.actionCommand == 'look') &&
      actionRec.actionTargetLocation != null) {
    direction = actionRec.actionTargetLocation?.locationDirection;
    log.fine(
      'Play command >${actionRec.actionCommand}< direction >$direction<',
    );
  }

  if (actionRec.actionTargetCharacter != null) {
    var target = actionRec.actionTargetCharacter;
    log.fine(
      'Play command >${actionRec.actionCommand}< character >${target?.characterName}<',
    );
  }

  if (actionRec.actionTargetMonster != null) {
    var target = actionRec.actionTargetMonster;
    log.fine(
      'Play command >${actionRec.actionCommand}< monster >${target?.monsterName}<',
    );
  }

  if (actionRec.actionTargetObject != null) {
    var target = actionRec.actionTargetObject;
    log.fine(
      'Play command >${actionRec.actionCommand}< object >${target?.objectName}<',
    );
  }
}
