import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';

class GameStashedWidget extends StatefulWidget {
  const GameStashedWidget({Key? key}) : super(key: key);

  @override
  _GameStashedWidgetState createState() => _GameStashedWidgetState();
}

class _GameStashedWidgetState extends State<GameStashedWidget> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameStashedWidget');
    log.info('Building..');

    return BlocConsumer<DungeonActionCubit, DungeonActionState>(
      listener: (BuildContext context, DungeonActionState state) {
        log.info('listener...');
      },
      builder: (BuildContext context, DungeonActionState state) {
        log.info('Rendering stashed inventory');

        return const Text('Stashed');
      },
    );
  }
}
