import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon_character/dungeon_character_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_command/dungeon_command_cubit.dart';

void submitLookAction(BuildContext context) async {
  final log = getLogger('GameWidget', '_initAction');
  log.fine('Initialising action..');

  final dungeonCharacterCubit = BlocProvider.of<DungeonCharacterCubit>(context);
  if (dungeonCharacterCubit.dungeonCharacterRecord == null) {
    log.warning(
        'Dungeon character cubit missing dungeon character record, cannot initialise action');
    return;
  }

  final dungeonActionCubit = BlocProvider.of<DungeonActionCubit>(context);

  return dungeonActionCubit
      .createAction(
        dungeonCharacterCubit.dungeonCharacterRecord!.dungeonID,
        dungeonCharacterCubit.dungeonCharacterRecord!.characterID,
        "look",
      )
      .then((v) => playActions(context));
}

Future<void> submitAction(BuildContext context) async {
  final log = getLogger('Action', 'submitAction');

  final dungeonCharacterCubit = BlocProvider.of<DungeonCharacterCubit>(context);
  if (dungeonCharacterCubit.dungeonCharacterRecord == null) {
    log.warning(
        'Dungeon character cubit missing dungeon character record, cannot initialise action');
    return;
  }

  final dungeonActionCubit = BlocProvider.of<DungeonActionCubit>(context);
  final dungeonCommandCubit = BlocProvider.of<DungeonCommandCubit>(context);

  log.warning('Submitting action ${dungeonCommandCubit.command()}');

  return dungeonActionCubit
      .createAction(
        dungeonCharacterCubit.dungeonCharacterRecord!.dungeonID,
        dungeonCharacterCubit.dungeonCharacterRecord!.characterID,
        dungeonCommandCubit.command(),
      )
      .then((v) => playActions(context));
}

// playActions plays all character and other available actions
Future<void> playActions(BuildContext context) async {
  final log = getLogger('Action', 'playAction');

  final dungeonActionCubit = BlocProvider.of<DungeonActionCubit>(context);

  log.fine('Playing character action');
  dungeonActionCubit.playCharacterAction();

  log.fine('Playing other actions');
  dungeonActionCubit.playOtherActions();
}
