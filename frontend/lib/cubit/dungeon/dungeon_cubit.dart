import 'package:bloc/bloc.dart';
import 'package:meta/meta.dart';
import 'package:equatable/equatable.dart';

// Application
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/repository/repository.dart';

part 'dungeon_state.dart';

class DungeonCubit extends Cubit<DungeonState> {
  final Map<String, String> config;
  final RepositoryCollection repositories;

  List<DungeonRecord>? dungeonRecords;
  DungeonRecord? dungeonRecord;

  DungeonCubit({required this.config, required this.repositories})
      : super(const DungeonStateInitial());

  // TODO: Unused
  void clearDungeon() {
    dungeonRecord = null;
    emit(DungeonStateLoaded(dungeonRecords: dungeonRecords));
  }

  Future<void> loadDungeons() async {
    final log = getLogger('DungeonCubit');
    log.fine('Loading dungeons...');
    emit(const DungeonStateLoading());

    List<DungeonRecord>? dungeonRecords;

    try {
      dungeonRecords = await repositories.dungeonRepository.getMany();
    } catch (err) {
      emit(const DungeonStateLoadError());
      return;
    }

    emit(DungeonStateLoaded(dungeonRecords: dungeonRecords));
  }
}
