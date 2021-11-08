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

  Future<void> createAction(String dungeonID, String characterID, String command) async {
    final log = getLogger('DungeonActionCubit');
    log.info('Creating dungeon action command >$command<');

    emit(DungeonActionStateCreating(sentence: command));

    DungeonActionRecord? createdDungeonActionRecord =
        await repositories.dungeonActionRepository.create(dungeonID, characterID, command);

    log.info('Created dungeon action $createdDungeonActionRecord');

    if (createdDungeonActionRecord != null) {
      dungeonActionRecord = createdDungeonActionRecord;
      dungeonActionRecords.add(createdDungeonActionRecord);

      emit(
        DungeonActionStateCreated(dungeonActionRecord: createdDungeonActionRecord),
      );
    }
  }

  Future<bool> playAction(String dungeonID, String characterID, String command) async {
    final log = getLogger('DungeonActionCubit');
    log.info('Creating dungeon action command >$command<');

    if (dungeonActionRecords.length <= 1) {
      return false;
    }

    emit(
      DungeonActionStatePlaying(
          previous: dungeonActionRecords.removeAt(0), current: dungeonActionRecords[0]),
    );

    if (dungeonActionRecords.length <= 1) {
      return false;
    }

    return true;
  }
}
