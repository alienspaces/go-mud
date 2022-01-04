import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';

class GameInventoryWidget extends StatefulWidget {
  const GameInventoryWidget({Key? key}) : super(key: key);

  @override
  _GameInventoryWidgetState createState() => _GameInventoryWidgetState();
}

class _GameInventoryWidgetState extends State<GameInventoryWidget> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameInventoryWidget');
    log.info('Building..');

    return BlocConsumer<DungeonActionCubit, DungeonActionState>(
      listener: (BuildContext context, DungeonActionState state) {
        log.info('listener...');
      },
      builder: (BuildContext context, DungeonActionState state) {
        log.info('Rendering inventory');

        return const Text('Inventory');
      },
    );
  }
}
