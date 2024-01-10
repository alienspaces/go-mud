import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'dungeon_command/dungeon_command_cubit.dart';

void selectTarget(BuildContext context, String target) {
  final log = getLogger('Cubit', 'selectTarget');
  log.fine('Submitting move action..');

  final dungeonCommandCubit = BlocProvider.of<DungeonCommandCubit>(context);
  if (dungeonCommandCubit.target == target) {
    log.fine('++ Unselecting target $target');
    dungeonCommandCubit.unselectTarget();
    return;
  }

  log.fine('++ Selecting target $target');
  dungeonCommandCubit.selectTarget(target);
}
