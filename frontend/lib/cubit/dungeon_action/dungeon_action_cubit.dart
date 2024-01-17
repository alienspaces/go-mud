import 'package:bloc/bloc.dart';
import 'package:meta/meta.dart';
import 'package:equatable/equatable.dart';
import 'package:logging/logging.dart';
import 'dart:async';

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

  DungeonActionCubit({
    required this.config,
    required this.repositories,
  }) : super(const DungeonActionStateInitial());

  // Experiment to see what functionality the following methods might provide
  // to assist keeping cubit state syncronised and current.
  @override
  void onChange(Change<DungeonActionState> change) {
    super.onChange(change);
    final log = getLogger('DungeonActionCubit', 'onChange');
    log.fine("Changed:", change);
  }

  @override
  void onError(Object error, StackTrace stackTrace) {
    final log = getLogger('DungeonActionCubit', 'onError');
    log.fine("Errored:", error);
    super.onError(error, stackTrace);
  }
  // ^^^^ //

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
        type: ActionStateErrorType.tooEarly,
        message: 'Command too early, slow down...',
      ));
      return;
    } on ActionInvalidCharacterException {
      emit(const DungeonActionStateError(
        type: ActionStateErrorType.invalidCharacter,
        message: 'Character does not exist.',
      ));
      return;
    } on ActionInvalidDungeonException {
      emit(const DungeonActionStateError(
        type: ActionStateErrorType.invalidDungeon,
        message: 'Dungeon does not exist.',
      ));
      return;
    } on RepositoryException catch (err) {
      log.warning('Throwing action state error ${err.message}');
      emit(DungeonActionStateError(
        type: ActionStateErrorType.unknown,
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
      otherActionRecs.add(actionRec);
    }

    log.warning('**** Emitting DungeonActionStateCreated');

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

  /// Clears all action records
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
      log.fine('Other actions are empty, cannot play other action');
      return;
    }

    log.warning('Other actions length ${otherActionRecs.length}');

    var milliseconds = 350;
    while (otherActionRecs.isNotEmpty) {
      // Play from the earliest other action
      DungeonActionRecord actionRec = otherActionRecs.removeAt(0);

      String? actionDirection;
      if (actionRec.actionCommand == 'move' ||
          actionRec.actionCommand == 'look') {
        actionDirection = actionRec.actionTargetLocation?.locationDirection;
      }

      // Do not spam events, otherwise we probably need to look at streaming
      Timer(Duration(milliseconds: milliseconds), () {
        emit(
          DungeonActionStatePlayingOther(
            actionRec: actionRec,
            actionCommand: actionRec.actionCommand,
            actionDirection: actionDirection,
          ),
        );
      });

      milliseconds += 350;
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
      'Command >${actionRec.actionCommand}< direction >$direction<',
    );
  }

  if (actionRec.actionTargetCharacter != null) {
    var target = actionRec.actionTargetCharacter;
    log.fine(
      'Command >${actionRec.actionCommand}< character >${target?.characterName}<',
    );
  }

  if (actionRec.actionTargetMonster != null) {
    var target = actionRec.actionTargetMonster;
    log.fine(
      'Command >${actionRec.actionCommand}< monster >${target?.monsterName}<',
    );
  }

  if (actionRec.actionTargetObject != null) {
    var target = actionRec.actionTargetObject;
    log.fine(
      'Command >${actionRec.actionCommand}< object >${target?.objectName}<',
    );
  }
}
