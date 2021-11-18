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

  // TODO: Keep these up to date as we play through
  DungeonActionRecord? currentDungeonActionRecord;
  DungeonActionRecord? nextDungeonActionRecord;

  DungeonActionCubit({required this.config, required this.repositories})
      : super(const DungeonActionStateInitial());

  Future<void> createAction(String dungeonID, String characterID, String command) async {
    final log = getLogger('DungeonActionCubit');
    log.info('(createAction) Creating dungeon action command >$command<');

    emit(DungeonActionStateCreating(sentence: command));

    DungeonActionRecord? createdDungeonActionRecord =
        await repositories.dungeonActionRepository.create(dungeonID, characterID, command);

    log.info('(createAction) Created dungeon action $createdDungeonActionRecord');

    if (createdDungeonActionRecord != null) {
      dungeonActionRecord = createdDungeonActionRecord;
      dungeonActionRecords.add(createdDungeonActionRecord);

      emit(
        DungeonActionStateCreated(
          current: createdDungeonActionRecord,
          action: createdDungeonActionRecord.action.command,
          direction: createdDungeonActionRecord.action.targetDungeonLocationDirection,
        ),
      );
    }
  }

  /// Returns true when there are more actions to play
  bool playAction() {
    final log = getLogger('DungeonActionCubit');

    if (dungeonActionRecords.length < 2) {
      log.info('(playAction) Not enough dungeon action records, not playing action');
      return false;
    }

    DungeonActionRecord previous = dungeonActionRecords.removeAt(0);
    DungeonActionRecord current = dungeonActionRecords[0];
    String? direction;
    if (current.action.command == 'move' || current.action.command == 'look') {
      direction = current.action.targetDungeonLocationDirection;
    }

    log.info('(playAction) Play action command >${current.action.command}< direction >$direction<');

    emit(
      DungeonActionStatePlaying(
        previous: previous,
        current: current,
        action: current.action.command,
        direction: direction,
      ),
    );

    if (dungeonActionRecords.length <= 1) {
      return false;
    }

    return true;
  }
}
