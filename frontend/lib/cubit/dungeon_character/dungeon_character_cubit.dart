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

  /// getDungeonCharacterRecordForCharacter returns a DungeonCharacterRecord if
  /// the character is currently inside a dungeon
  Future<DungeonCharacterRecord?> getDungeonCharacterRecord(
      String dungeonID, String characterID) async {
    final log = getLogger('getDungeonCharacterRecordForCharacter');
    log.info('Searching through ${dungeonCharacterRecords?.length}');

    emit(DungeonCharacterStateLoading(characterID: characterID));

    // Find existing cached record
    DungeonCharacterRecord? rec = dungeonCharacterRecords?.firstWhere((rec) {
      log.info(
          'Testing dungeon ID ${rec.dungeonID} character ID ${rec.characterID}');
      if (rec.characterID == characterID) {
        return true;
      }
      return false;
    });

    if (rec != null) {
      log.info("Returning existing record");
      return rec;
    }

    try {
      rec = await repositories.dungeonCharacterRepository
          .getOne(dungeonID, characterID);
    } on RepositoryException catch (err) {
      log.warning('Throwing dungeon character load error');
      emit(DungeonCharacterStateLoadError(
          characterID: characterID, message: err.message));
      return Future<DungeonCharacterRecord?>.value(null);
    }

    return Future<DungeonCharacterRecord?>.value(rec);
  }

  Future<void> enterDungeonCharacter(
      String dungeonID, String characterID) async {
    final log = getLogger('DungeonCharacterCubit');
    log.info('Entering dungeon ID $dungeonID character ID $characterID');

    emit(const DungeonCharacterStateCreate());

    dungeonCharacterRecords?.forEach((rec) {
      log.info(
          'Existing dungeon ID ${rec.dungeonID} character ID ${rec.characterID}');
    });

    DungeonCharacterRecord? existingDungeonCharacterRecord =
        await getDungeonCharacterRecord(dungeonID, characterID);

    // Character already inside this dungeon
    if (existingDungeonCharacterRecord != null &&
        existingDungeonCharacterRecord.dungeonID == dungeonID) {
      log.info(
          'Dungeon with character $existingDungeonCharacterRecord is already in this dungeon');
      dungeonCharacterRecords?.add(existingDungeonCharacterRecord);
      dungeonCharacterRecord = existingDungeonCharacterRecord;
      emit(DungeonCharacterStateCreated(
          dungeonCharacterRecord: existingDungeonCharacterRecord));
      return;
    }

    // Character already inside some other dungeon
    if (existingDungeonCharacterRecord != null) {
      log.info(
          'Dungeon with character $existingDungeonCharacterRecord is already in a dungeon');
      emit(DungeonCharacterStateCreateError(
          dungeonID: dungeonID,
          characterID: characterID,
          message:
              'Dungeon with character $existingDungeonCharacterRecord is already in a dungeon'));
      return;
    }

    DungeonCharacterRecord? createdDungeonCharacterRecord;
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
      if (dungeonCharacterRecord.characterID == characterID &&
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
