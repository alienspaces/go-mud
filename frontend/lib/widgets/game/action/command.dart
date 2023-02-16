import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/widgets/game/action/action.dart';
import 'package:go_mud_client/cubit/dungeon_command/dungeon_command_cubit.dart';

class GameActionCommandWidget extends StatefulWidget {
  const GameActionCommandWidget({Key? key}) : super(key: key);

  @override
  State<GameActionCommandWidget> createState() =>
      _GameActionCommandWidgetState();
}

// TODO: 9-implement-monster-actions: Timer with animation that submits a
// look action an action is not performed by the player.

class _GameActionCommandWidgetState extends State<GameActionCommandWidget> {
  @override
  void initState() {
    final log = getLogger('GameWidget', 'initState');
    log.fine('Initialising state..');

    super.initState();

    submitCubitLookAction(context);
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameActionCommandWidget', 'build');
    log.fine('Building..');

    return BlocConsumer<DungeonCommandCubit, DungeonCommandState>(
      listener: (BuildContext context, DungeonCommandState state) {
        log.fine('listener...');
      },
      builder: (BuildContext context, DungeonCommandState state) {
        if (state is DungeonCommandStatePreparing) {
          // ignore: avoid_unnecessary_containers
          return Container(
            color: Colors.brown[200],
            alignment: Alignment.center,
            child:
                Text('${state.action ?? ''} ${state.target ?? ''}'.trimRight()),
          );
        }
        return Container();
      },
    );
  }
}
