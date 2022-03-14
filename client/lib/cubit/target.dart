import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';

// Local
import 'dungeon/dungeon_cubit.dart';
import 'dungeon_command/dungeon_command_cubit.dart';
import 'character/character_cubit.dart';

void selectTarget(BuildContext context, String target) {
  final log = getLogger('selectTarget');
  log.fine('Submitting move action..');

  final dungeonCubit = BlocProvider.of<DungeonCubit>(context);
  if (dungeonCubit.dungeonRecord == null) {
    log.warning(
        'Dungeon cubit missing dungeon record, cannot initialise action');
    return;
  }

  final characterCubit = BlocProvider.of<CharacterCubit>(context);
  if (characterCubit.characterRecord == null) {
    log.warning(
        'Character cubit missing character record, cannot initialise action');
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
