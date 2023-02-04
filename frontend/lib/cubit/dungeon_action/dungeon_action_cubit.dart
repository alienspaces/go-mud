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

    DungeonActionRecord? createdDungeonActionRecord;
    try {
      createdDungeonActionRecord = await repositories.dungeonActionRepository
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

    log.info(
        'location ${createdDungeonActionRecord?.actionLocation.locationName}');

    if (createdDungeonActionRecord?.actionTargetLocation != null) {
      log.info(
          'targetLocation ${createdDungeonActionRecord?.actionTargetLocation?.locationName}');
    }
    if (createdDungeonActionRecord?.actionTargetCharacter != null) {
      log.info(
          'targetCharacter ${createdDungeonActionRecord?.actionTargetCharacter?.toJson()}');
    }
    if (createdDungeonActionRecord?.actionTargetMonster != null) {
      log.info(
          'targetMonster ${createdDungeonActionRecord?.actionTargetMonster?.toJson()}');
    }
    if (createdDungeonActionRecord?.actionTargetObject != null) {
      log.info(
          'targetObject ${createdDungeonActionRecord?.actionTargetObject?.toJson()}');
    }

    if (createdDungeonActionRecord != null) {
      dungeonActionRecord = createdDungeonActionRecord;
      dungeonActionRecords.add(createdDungeonActionRecord);

      emit(
        DungeonActionStateCreated(
          current: createdDungeonActionRecord,
          action: createdDungeonActionRecord.actionCommand,
          direction: createdDungeonActionRecord
              .actionTargetLocation?.locationDirection,
        ),
      );
    }
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
      log.fine('Not enough dungeon action records, not playing action');
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
