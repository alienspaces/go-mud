import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';

import 'package:go_mud_client/cubit/dungeon_character/dungeon_character_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_command/dungeon_command_cubit.dart';

void selectAction(BuildContext context, String action) {
  final log = getLogger('GameActionPanelWidget');
  log.fine('Selecting action..');

  final dungeonCharacterCubit = BlocProvider.of<DungeonCharacterCubit>(context);
  if (dungeonCharacterCubit.dungeonCharacterRecord == null) {
    log.warning(
        'Dungeon character cubit missing dungeon character record, cannot initialise action');
    return;
  }

  final dungeonCommandCubit = BlocProvider.of<DungeonCommandCubit>(context);
  if (dungeonCommandCubit.action == action) {
    log.fine('++ Unselecting action $action');
    dungeonCommandCubit.unselectAction();
    return;
  }

  log.fine('++ Selecting action $action');
  dungeonCommandCubit.selectAction(action);
}

void submitAction(BuildContext context) async {
  final log = getLogger('GameActionPanelWidget');
  log.fine('Submitting action..');

  final dungeonCharacterCubit = BlocProvider.of<DungeonCharacterCubit>(context);
  if (dungeonCharacterCubit.dungeonCharacterRecord == null) {
    log.warning(
        'Dungeon character cubit missing dungeon character record, cannot initialise action');
    return;
  }

  log.fine('++ Submitting action');
  final dungeonActionCubit = BlocProvider.of<DungeonActionCubit>(context);
  final dungeonCommandCubit = BlocProvider.of<DungeonCommandCubit>(context);

  await dungeonActionCubit.createAction(
    dungeonCharacterCubit.dungeonCharacterRecord!.dungeonID,
    dungeonCharacterCubit.dungeonCharacterRecord!.characterID,
    dungeonCommandCubit.command(),
  );
  dungeonCommandCubit.unselectAll();

  // TODO: (client) Loop this using a timer allowing animations to complete
  var moreActions = dungeonActionCubit.playAction();
  log.fine('++ More actions >$moreActions<');
}
