import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';

// Local
import 'dungeon_character/dungeon_character_cubit.dart';
import 'dungeon_command/dungeon_command_cubit.dart';

void selectTarget(BuildContext context, String target) {
  final log = getLogger('selectTarget');
  log.fine('Submitting move action..');

  final dungeonCharacterCubit = BlocProvider.of<DungeonCharacterCubit>(context);
  if (dungeonCharacterCubit.dungeonCharacterRecord == null) {
    log.warning(
        'Dungeon character cubit missing dungeon character record, cannot initialise action');
    return;
  }

  final dungeonCommandCubit = BlocProvider.of<DungeonCommandCubit>(context);
  if (dungeonCommandCubit.target == target) {
    log.fine('++ Unselecting target $target');
    dungeonCommandCubit.unselectTarget();
    return;
  }

  log.fine('++ Selecting target $target');
  dungeonCommandCubit.selectTarget(target);
}
