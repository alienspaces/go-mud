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
    final log = getLogger('DungeonActionCubit');
    log.fine('(createAction) Creating dungeon action command >$command<');

    emit(DungeonActionStateCreating(
      sentence: command,
      current: dungeonActionRecord,
    ));

    DungeonActionRecord? createdDungeonActionRecord = await repositories
        .dungeonActionRepository
        .create(dungeonID, characterID, command);

    log.fine(
        '(createAction) location ${createdDungeonActionRecord?.actionLocation.locationName}');
    log.fine(
        '(createAction) targetLocation ${createdDungeonActionRecord?.actionTargetLocation?.locationName}');
    log.fine(
        '(createAction) targetCharacter ${createdDungeonActionRecord?.actionTargetCharacter?.characterName}');
    log.fine(
        '(createAction) targetMonster ${createdDungeonActionRecord?.actionTargetMonster?.monsterName}');
    log.fine(
        '(createAction) targetObject ${createdDungeonActionRecord?.actionTargetObject?.objectName}');

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
    final log = getLogger('DungeonActionCubit');

    if (dungeonActionRecords.length < 2) {
      log.fine(
          '(playAction) Not enough dungeon action records, not playing action');
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
          '(playAction) Play action command >${current.actionCommand}< direction >$direction<');
    }

    if (current.actionTargetCharacter != null) {
      log.fine(
          '(playAction) Play action command >${current.actionCommand}< character >${current.actionTargetCharacter?.characterName}<');
    }

    if (current.actionTargetMonster != null) {
      log.fine(
          '(playAction) Play action command >${current.actionCommand}< monster >${current.actionTargetMonster?.monsterName}<');
    }

    if (current.actionTargetObject != null) {
      log.fine(
          '(playAction) Play action command >${current.actionCommand}< object >${current.actionTargetObject?.objectName}<');
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
