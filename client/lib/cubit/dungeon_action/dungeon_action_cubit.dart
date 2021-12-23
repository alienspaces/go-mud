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
    log.info('(createAction) Creating dungeon action command >$command<');

    emit(DungeonActionStateCreating(
      sentence: command,
      current: dungeonActionRecord,
    ));

    DungeonActionRecord? createdDungeonActionRecord = await repositories
        .dungeonActionRepository
        .create(dungeonID, characterID, command);

    log.info(
        '(createAction) location ${createdDungeonActionRecord?.location.name}');
    log.info(
        '(createAction) targetLocation ${createdDungeonActionRecord?.targetLocation?.name}');
    log.info(
        '(createAction) targetCharacter ${createdDungeonActionRecord?.targetCharacter?.name}');
    log.info(
        '(createAction) targetMonster ${createdDungeonActionRecord?.targetMonster?.name}');
    log.info(
        '(createAction) targetObject ${createdDungeonActionRecord?.targetObject?.name}');

    if (createdDungeonActionRecord != null) {
      dungeonActionRecord = createdDungeonActionRecord;
      dungeonActionRecords.add(createdDungeonActionRecord);

      emit(
        DungeonActionStateCreated(
          current: createdDungeonActionRecord,
          action: createdDungeonActionRecord.command,
          direction: createdDungeonActionRecord.targetLocation?.direction,
        ),
      );
    }
  }

  /// Returns true when there are more actions to play
  bool playAction() {
    final log = getLogger('DungeonActionCubit');

    if (dungeonActionRecords.length < 2) {
      log.info(
          '(playAction) Not enough dungeon action records, not playing action');
      return false;
    }

    DungeonActionRecord previous = dungeonActionRecords.removeAt(0);
    DungeonActionRecord current = dungeonActionRecords[0];
    String? direction;
    if (current.command == 'move' || current.command == 'look') {
      direction = current.targetLocation?.direction;
    }

    if (current.targetLocation != null) {
      log.info(
          '(playAction) Play action command >${current.command}< direction >$direction<');
    }

    if (current.targetCharacter != null) {
      log.info(
          '(playAction) Play action command >${current.command}< character >${current.targetCharacter?.name}<');
    }

    if (current.targetMonster != null) {
      log.info(
          '(playAction) Play action command >${current.command}< monster >${current.targetMonster?.name}<');
    }

    if (current.targetObject != null) {
      log.info(
          '(playAction) Play action command >${current.command}< object >${current.targetObject?.name}<');
    }

    emit(
      DungeonActionStatePlaying(
        previous: previous,
        current: current,
        action: current.command,
        direction: direction,
      ),
    );

    if (dungeonActionRecords.length <= 1) {
      return false;
    }

    return true;
  }
}
