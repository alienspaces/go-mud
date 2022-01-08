import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';

class GameBoardObjectWidget extends StatefulWidget {
  const GameBoardObjectWidget({Key? key}) : super(key: key);

  @override
  _GameBoardObjectWidgetState createState() => _GameBoardObjectWidgetState();
}

class _GameBoardObjectWidgetState extends State<GameBoardObjectWidget> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameBoardObjectWidget');
    log.info('Building..');

    return BlocConsumer<DungeonActionCubit, DungeonActionState>(
        listener: (BuildContext context, DungeonActionState state) {
      log.info('listener...');
    }, builder: (BuildContext context, DungeonActionState state) {
      if (state is DungeonActionStatePlaying) {
        var dungeonActionRecord = state.current;

        log.info(
            'DungeonActionStatePlaying - Rendering action ${dungeonActionRecord.command}');

        if (dungeonActionRecord.targetObject != null) {
          log.info('Rendering look target object');
          return Container(
            padding: const EdgeInsets.all(5),
            child: const Text("Looking object"),
          );
        }
      }

      return Container();
    });
  }
}
