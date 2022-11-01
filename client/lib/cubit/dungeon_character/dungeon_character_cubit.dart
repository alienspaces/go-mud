import 'package:bloc/bloc.dart';
import 'package:meta/meta.dart';
import 'package:equatable/equatable.dart';

// Application
import 'package:go_mud_client/repository/repository.dart';
import 'package:go_mud_client/logger.dart';

part 'dungeon_character_state.dart';

class DungeonCharacterCubit extends Cubit<DungeonCharacterState> {
  final Map<String, String> config;
  final RepositoryCollection repositories;

  List<DungeonCharacterRecord>? dungeonCharacterRecords;
  DungeonCharacterRecord? dungeonCharacterRecord;

  DungeonCharacterCubit({required this.config, required this.repositories})
      : super(const DungeonCharacterStateInitial());

  Future<void> enterDungeonCharacter(
      String dungeonID, String characterID) async {
    final log = getLogger('DungeonCharacterCubit');
    log.info('Entering dungeon ID $dungeonID character ID $characterID');

    emit(const DungeonCharacterStateCreate());

    DungeonCharacterRecord? createdDungeonCharacterRecord;

    dungeonCharacterRecords?.forEach((rec) {
      log.info('Existing dungeon ID ${rec.dungeonID} character ID ${rec.id}');
    });

    createdDungeonCharacterRecord = dungeonCharacterRecords?.firstWhere((rec) {
      log.info('Testing dungeon ID ${rec.dungeonID} character ID ${rec.id}');
      if (rec.id == characterID && rec.dungeonID == dungeonID) {
        return true;
      }
      return false;
    });

    if (createdDungeonCharacterRecord != null) {
      log.info(
          'Dungeon with character $createdDungeonCharacterRecord is already in a dungeon');
    }

    if (createdDungeonCharacterRecord == null) {
      try {
        createdDungeonCharacterRecord = await repositories
            .dungeonCharacterRepository
            .createOne(dungeonID, characterID);
      } on RepositoryException catch (err) {
        log.warning('Throwing dungeon character enter error');
        emit(DungeonCharacterStateCreateError(
            dungeonID: dungeonID,
            characterID: characterID,
            message: err.message));
        return;
      }
    }

    if (createdDungeonCharacterRecord != null) {
      log.info('Entered dungeon with character $createdDungeonCharacterRecord');
      dungeonCharacterRecords?.add(createdDungeonCharacterRecord);
      dungeonCharacterRecord = createdDungeonCharacterRecord;
      emit(DungeonCharacterStateCreated(
          dungeonCharacterRecord: createdDungeonCharacterRecord));
    }
  }

  Future<void> exitDungeonCharacter(
      String dungeonID, String characterID) async {
    final log = getLogger('DungeonCharacterCubit');
    log.fine('Exiting dungeon ID $dungeonID character ID $characterID');

    emit(const DungeonCharacterStateDelete());

    DungeonCharacterRecord? deletedDungeonCharacterRecord =
        dungeonCharacterRecords?.firstWhere((dungeonCharacterRecord) {
      if (dungeonCharacterRecord.id == characterID &&
          dungeonCharacterRecord.dungeonID == dungeonID) {
        return true;
      }
      return false;
    });

    if (deletedDungeonCharacterRecord == null) {
      emit(DungeonCharacterStateDeleteError(
          dungeonID: dungeonID,
          characterID: characterID,
          message: 'Exiting dungeon character record does not exist'));
      return;
    }

    try {
      await repositories.dungeonCharacterRepository
          .deleteOne(dungeonID, characterID);
    } on RepositoryException catch (err) {
      log.warning('Throwing dungeon character exit error');
      emit(DungeonCharacterStateDeleteError(
          dungeonID: dungeonID,
          characterID: characterID,
          message: err.message));
      return;
    }

    log.info('Exited dungeon character $deletedDungeonCharacterRecord');
    dungeonCharacterRecords?.remove(deletedDungeonCharacterRecord);
    emit(DungeonCharacterStateDeleted(
        dungeonCharacterRecord: deletedDungeonCharacterRecord));
  }
}
