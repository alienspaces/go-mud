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

  DungeonCharacterRecord? dungeonCharacterRecord;

  DungeonCharacterCubit({required this.config, required this.repositories})
      : super(const DungeonCharacterStateInitial());

  /// getDungeonCharacterRecordForCharacter returns a DungeonCharacterRecord if
  /// the character is currently inside a dungeon
  Future<DungeonCharacterRecord?> getDungeonCharacterRecord(
      String dungeonID, String characterID) async {
    final log = getLogger('getDungeonCharacterRecordForCharacter');

    emit(DungeonCharacterStateLoading(characterID: characterID));

    DungeonCharacterRecord? rec;
    try {
      rec = await repositories.dungeonCharacterRepository
          .getDungeonCharacter(dungeonID, characterID);
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

    dungeonCharacterRecord =
        await getDungeonCharacterRecord(dungeonID, characterID);

    // Character already inside some other dungeon
    if (dungeonCharacterRecord != null) {
      log.info(
          'Dungeon with character $dungeonCharacterRecord is already in a dungeon');
      emit(DungeonCharacterStateCreateError(
          dungeonID: dungeonID,
          characterID: characterID,
          message:
              'Dungeon with character $dungeonCharacterRecord is already in a dungeon'));
      return;
    }

    try {
      dungeonCharacterRecord = await repositories.dungeonCharacterRepository
          .enterDungeonCharacter(dungeonID, characterID);
    } on RepositoryException catch (err) {
      log.warning('Throwing dungeon character enter error');
      emit(DungeonCharacterStateCreateError(
          dungeonID: dungeonID,
          characterID: characterID,
          message: err.message));
      return;
    }

    if (dungeonCharacterRecord != null) {
      log.info('Entered dungeon with character $dungeonCharacterRecord');
      emit(DungeonCharacterStateCreated(
          dungeonCharacterRecord: dungeonCharacterRecord!));
    }
  }

  Future<void> exitDungeonCharacter(
      String dungeonID, String characterID) async {
    final log = getLogger('DungeonCharacterCubit');
    log.info('Exiting dungeon ID $dungeonID character ID $characterID');

    emit(const DungeonCharacterStateDelete());

    DungeonCharacterRecord? dungeonCharacterRecord =
        await getDungeonCharacterRecord(dungeonID, characterID);

    if (dungeonCharacterRecord == null) {
      log.warning(
          'Did not find dungeon character record for dungeon ID $dungeonID character ID $characterID');

      emit(DungeonCharacterStateDeleteError(
          dungeonID: dungeonID,
          characterID: characterID,
          message: 'Exiting dungeon character record does not exist'));
      return;
    }

    try {
      await repositories.dungeonCharacterRepository
          .exitDungeonCharacter(dungeonID, characterID);
    } on RepositoryException catch (err) {
      log.warning('Throwing dungeon character exit error');
      emit(DungeonCharacterStateDeleteError(
          dungeonID: dungeonID,
          characterID: characterID,
          message: err.message));
      return;
    }

    log.info('Exited dungeon character $dungeonCharacterRecord');
    emit(DungeonCharacterStateDeleted(
        dungeonCharacterRecord: dungeonCharacterRecord));
    dungeonCharacterRecord = null;
  }
}
