import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';

class GameBoardMonsterWidget extends StatefulWidget {
  const GameBoardMonsterWidget({Key? key}) : super(key: key);

  @override
  _GameBoardMonsterWidgetState createState() => _GameBoardMonsterWidgetState();
}

class _GameBoardMonsterWidgetState extends State<GameBoardMonsterWidget> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameBoardMonsterWidget');
    log.info('Building..');

    return BlocConsumer<DungeonActionCubit, DungeonActionState>(
        listener: (BuildContext context, DungeonActionState state) {
      log.info('listener...');
    }, builder: (BuildContext context, DungeonActionState state) {
      if (state is DungeonActionStatePlaying) {
        var dungeonActionRecord = state.current;

        log.info(
            'DungeonActionStatePlaying - Rendering action ${dungeonActionRecord.command}');

        if (dungeonActionRecord.targetMonster != null) {
          log.info('Rendering look target monster');
          return Container(
            padding: const EdgeInsets.all(5),
            child: const Text("Looking monster"),
          );
        }
      }

      return Container();
    });
  }
}
