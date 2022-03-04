import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon/dungeon_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_command/dungeon_command_cubit.dart';
import 'package:go_mud_client/cubit/character/character_cubit.dart';

// TODO: USe from button widgets
void selectTarget(BuildContext context, String target) {
  final log = getLogger('selectTarget');
  log.info('Submitting move action..');

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
    log.info('++ Unselecting target $target');
    dungeonCommandCubit.unselectTarget();
    return;
  }

  log.info('++ Selecting target $target');
  dungeonCommandCubit.selectTarget(target);
}
