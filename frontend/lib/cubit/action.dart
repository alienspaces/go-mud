import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';

import 'package:go_mud_client/cubit/character/character_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_command/dungeon_command_cubit.dart';

void selectAction(BuildContext context, String action) {
  final log = getLogger('Cubit', 'selectAction');
  log.fine('Selecting action..');

  final characterCubit = BlocProvider.of<CharacterCubit>(context);

  var characterRecord = characterCubit.characterRecord;

  if (characterRecord == null || characterRecord.dungeonID != null) {
    log.warning(
        'Character cubit missing character record, cannot select action');
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
  final log = getLogger('Cubit', 'submitAction');
  log.fine('Submitting action..');

  final characterCubit = BlocProvider.of<CharacterCubit>(context);

  var characterRecord = characterCubit.characterRecord;

  if (characterRecord == null || characterRecord.dungeonID != null) {
    log.warning(
        'Character cubit missing character record, cannot submit action');
    return;
  }

  log.fine('++ Submitting action');
  final dungeonActionCubit = BlocProvider.of<DungeonActionCubit>(context);
  final dungeonCommandCubit = BlocProvider.of<DungeonCommandCubit>(context);

  await dungeonActionCubit.createAction(
    characterRecord.dungeonID!,
    characterRecord.characterID,
    dungeonCommandCubit.command(),
  );

  dungeonCommandCubit.unselectAll();
}
