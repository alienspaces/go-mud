import 'package:bloc/bloc.dart';
import 'package:meta/meta.dart';
import 'package:equatable/equatable.dart';

// Application
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/repository/repository.dart';

part 'dungeon_command_state.dart';

class DungeonCommandCubit extends Cubit<DungeonCommandState> {
  final Map<String, String> config;
  final RepositoryCollection repositories;

  String? action;
  String? target;

  DungeonCommandCubit({required this.config, required this.repositories})
      : super(const DungeonCommandStateInitial());

  String command() {
    return '$action $target'.trimRight();
  }

  void unselectAll() {
    final log = getLogger('DungeonCommandCubit');
    log.info('Unselecting action >$action< and target >$target<');
    action = null;
    target = null;
    emit(DungeonCommandStatePreparing(action: action, target: target));
  }

  Future<void> selectAction(String selectAction) async {
    final log = getLogger('DungeonCommandCubit');
    log.info('Selecting action >$selectAction<');
    action = selectAction;
    emit(DungeonCommandStatePreparing(action: action, target: target));
  }

  Future<void> unselectAction() async {
    final log = getLogger('DungeonCommandCubit');
    log.info('Unselecting action >$action<');
    action = null;
    target = null;
    emit(DungeonCommandStatePreparing(action: action, target: target));
  }

  Future<void> selectTarget(String selectTarget) async {
    final log = getLogger('DungeonCommandCubit');
    log.info('Selecting target >$selectTarget<');
    target = selectTarget;
    emit(DungeonCommandStatePreparing(action: action, target: target));
  }

  Future<void> unselectTarget() async {
    final log = getLogger('DungeonCommandCubit');
    log.info('Unselecting target >$target<');
    target = null;
    emit(DungeonCommandStatePreparing(action: action, target: target));
  }
}
