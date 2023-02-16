import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon_character/dungeon_character_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_command/dungeon_command_cubit.dart';

void submitCubitLookAction(BuildContext context) async {
  final log = getLogger('GameWidget', '_initAction');
  log.fine('Initialising action..');

  final dungeonCommandCubit = BlocProvider.of<DungeonCommandCubit>(context);
  dungeonCommandCubit.selectAction(
    'look',
  );

  await submitCubitAction(context);

  dungeonCommandCubit.unselectAction();
}

Future<void> submitCubitAction(BuildContext context) async {
  final log = getLogger('Action', 'submitAction');
  log.fine('Submitting action..');

  final dungeonCharacterCubit = BlocProvider.of<DungeonCharacterCubit>(context);
  if (dungeonCharacterCubit.dungeonCharacterRecord == null) {
    log.warning(
        'Dungeon character cubit missing dungeon character record, cannot initialise action');
    return;
  }

  final dungeonActionCubit = BlocProvider.of<DungeonActionCubit>(context);
  final dungeonCommandCubit = BlocProvider.of<DungeonCommandCubit>(context);

  return dungeonActionCubit.createAction(
    dungeonCharacterCubit.dungeonCharacterRecord!.dungeonID,
    dungeonCharacterCubit.dungeonCharacterRecord!.characterID,
    dungeonCommandCubit.command(),
  );
}
