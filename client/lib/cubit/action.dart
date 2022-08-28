import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';

import 'package:go_mud_client/cubit/dungeon/dungeon_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_command/dungeon_command_cubit.dart';
import 'package:go_mud_client/cubit/character/character_cubit.dart';

void selectAction(BuildContext context, String action) {
  final log = getLogger('GameActionPanelWidget');
  log.fine('Selecting action..');

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

  log.fine('++ Submitting action');
  final dungeonActionCubit = BlocProvider.of<DungeonActionCubit>(context);
  final dungeonCommandCubit = BlocProvider.of<DungeonCommandCubit>(context);

  await dungeonActionCubit.createAction(
    dungeonCubit.dungeonRecord!.id,
    characterCubit.characterRecord!.id,
    dungeonCommandCubit.command(),
  );
  dungeonCommandCubit.unselectAll();

  // TODO: (client) Loop this using a timer allowing animations to complete
  var moreActions = dungeonActionCubit.playAction();
  log.fine('++ More actions >$moreActions<');
}
