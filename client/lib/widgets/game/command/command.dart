import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon_command/dungeon_command_cubit.dart';

class GameCommandWidget extends StatefulWidget {
  const GameCommandWidget({Key? key}) : super(key: key);

  @override
  _GameCommandWidgetState createState() => _GameCommandWidgetState();
}

class _GameCommandWidgetState extends State<GameCommandWidget> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameLocationGridWidget');
    log.info('Building..');

    return BlocConsumer<DungeonCommandCubit, DungeonCommandState>(
      listener: (BuildContext context, DungeonCommandState state) {
        log.info('listener...');
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
