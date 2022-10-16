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

  DungeonCharacterCubit({required this.config, required this.repositories})
      : super(const DungeonCharacterStateInitial());

  Future<void> createDungeonCharacter(
      String dungeonID, String characterID) async {
    final log = getLogger('DungeonCharacterCubit');
    log.fine('Creating dungeon ID $dungeonID character ID $characterID');

    emit(const DungeonCharacterStateCreate());

    DungeonCharacterRecord? createdDungeonCharacterRecord;

    try {
      createdDungeonCharacterRecord = await repositories
          .dungeonCharacterRepository
          .createOne(dungeonID, characterID);
    } on RepositoryException catch (err) {
      log.warning('Throwing character create error');
      emit(DungeonCharacterStateCreateError(
          dungeonID: dungeonID,
          characterID: characterID,
          message: err.message));
      return;
    }

    if (createdDungeonCharacterRecord != null) {
      log.fine('Created dungeon character $createdDungeonCharacterRecord');
      dungeonCharacterRecords?.add(createdDungeonCharacterRecord);
      emit(DungeonCharacterStateCreated(
          dungeonCharacterRecord: createdDungeonCharacterRecord));
    }
  }

  Future<void> deleteDungeonCharacter(
      String dungeonID, String characterID) async {
    final log = getLogger('DungeonCharacterCubit');
    log.fine('Deleting dungeon ID $dungeonID character ID $characterID');

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
          message: 'Dungeon character record does not exist'));
      return;
    }

    try {
      await repositories.dungeonCharacterRepository
          .deleteOne(dungeonID, characterID);
    } on RepositoryException catch (err) {
      log.warning('Throwing character create error');
      emit(DungeonCharacterStateDeleteError(
          dungeonID: dungeonID,
          characterID: characterID,
          message: err.message));
      return;
    }

    log.fine('Deleted dungeon character $deletedDungeonCharacterRecord');
    dungeonCharacterRecords?.remove(deletedDungeonCharacterRecord);
    emit(DungeonCharacterStateDeleted(
        dungeonCharacterRecord: deletedDungeonCharacterRecord));
  }
}
