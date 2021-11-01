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

  String? action;
  String? target;

  List<DungeonActionRecord>? dungeonActionRecords;
  DungeonActionRecord? dungeonActionRecord;

  DungeonActionCubit({required this.config, required this.repositories})
      : super(const DungeonActionStateInitial());

  Future<void> selectAction(String selectAction) async {
    final log = getLogger('DungeonActionCubit');
    log.info('Selecting action $selectAction');
    action = selectAction;
    emit(DungeonActionStateCommand(action: action, target: target));
  }

  Future<void> unselectAction() async {
    final log = getLogger('DungeonActionCubit');
    log.info('Unselecting action $action');
    action = null;
    emit(DungeonActionStateCommand(action: action, target: target));
  }

  Future<void> selectTarget(String selectTarget) async {
    final log = getLogger('DungeonActionCubit');
    log.info('Selecting target $selectTarget');
    target = selectTarget;
    emit(DungeonActionStateCommand(action: action, target: target));
  }

  Future<void> unselectTarget() async {
    final log = getLogger('DungeonActionCubit');
    log.info('Unselecting target $target');
    target = null;
    emit(DungeonActionStateCommand(action: action, target: target));
  }

  Future<void> submitAction(String dungeonID, String characterID) async {
    final log = getLogger('DungeonActionCubit');
    log.info('Creating dungeon action >$action< target >$target<');

    emit(DungeonActionStateCreating(sentence: '$action $target'));

    DungeonActionRecord? createdDungeonActionRecord = await repositories.dungeonActionRepository
        .create(dungeonID, characterID, '$action $target');

    log.info('Created dungeon action $createdDungeonActionRecord');

    if (createdDungeonActionRecord != null) {
      dungeonActionRecord = createdDungeonActionRecord;
      emit(DungeonActionStateCreated(dungeonActionRecord: createdDungeonActionRecord));
    }
  }
}
