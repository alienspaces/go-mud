import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';

class GameBoardCharacterWidget extends StatefulWidget {
  const GameBoardCharacterWidget({Key? key}) : super(key: key);

  @override
  _GameBoardCharacterWidgetState createState() =>
      _GameBoardCharacterWidgetState();
}

class _GameBoardCharacterWidgetState extends State<GameBoardCharacterWidget> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameBoardCharacterWidget');
    log.info('Building..');

    return BlocConsumer<DungeonActionCubit, DungeonActionState>(
        listener: (BuildContext context, DungeonActionState state) {
      log.info('listener...');
    }, builder: (BuildContext context, DungeonActionState state) {
      if (state is DungeonActionStatePlaying) {
        var dungeonActionRecord = state.current;

        log.info(
            'DungeonActionStatePlaying - Rendering action ${dungeonActionRecord.command}');

        if (dungeonActionRecord.targetCharacter != null) {
          log.info('Rendering look target character');
          return Container(
            padding: const EdgeInsets.all(5),
            child: const Text("Looking character"),
          );
        }
      }

      return Container();
    });
  }
}
