import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';

class GameInventoryEquippedWidget extends StatefulWidget {
  const GameInventoryEquippedWidget({Key? key}) : super(key: key);

  @override
  _GameInventoryEquippedWidgetState createState() =>
      _GameInventoryEquippedWidgetState();
}

class _GameInventoryEquippedWidgetState
    extends State<GameInventoryEquippedWidget> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameInventoryEquippedWidget');
    log.info('Building..');

    return BlocConsumer<DungeonActionCubit, DungeonActionState>(
      listener: (BuildContext context, DungeonActionState state) {
        log.info('listener...');
      },
      builder: (BuildContext context, DungeonActionState state) {
        log.info('Rendering equipped inventory');

        return const Text('Equipped');
      },
    );
  }
}
