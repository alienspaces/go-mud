import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon_character/dungeon_character_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_command/dungeon_command_cubit.dart';

// initLoop submit the initial action and then monitors for additional
// actions submitted by the player. When the player fails to submit an
// action within a configurable number of seconds the default action of
// "look" will be automatically submitted. This could potentially be a
// part of the command widget with a visible animated timer as well.
void initLoop(BuildContext context) {
  final log = getLogger('GameWidget', '_initAction');
  log.fine('Initialising action..');

  final dungeonCharacterCubit = BlocProvider.of<DungeonCharacterCubit>(context);
  if (dungeonCharacterCubit.dungeonCharacterRecord == null) {
    log.warning(
        'Dungeon character cubit missing dungeon character record, cannot initialise action');
    return;
  }

  final dungeonCommandCubit = BlocProvider.of<DungeonCommandCubit>(context);
  dungeonCommandCubit.selectAction(
    'look',
  );

  final dungeonActionCubit = BlocProvider.of<DungeonActionCubit>(context);
  dungeonActionCubit.createAction(
    dungeonCharacterCubit.dungeonCharacterRecord!.dungeonID,
    dungeonCharacterCubit.dungeonCharacterRecord!.characterID,
    dungeonCommandCubit.command(),
  );

  dungeonCommandCubit.unselectAction();
}
