import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';

import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';

import 'package:go_mud_client/widgets/game/board/buttons/object_button.dart';

class GameEquippedWidget extends StatefulWidget {
  const GameEquippedWidget({Key? key}) : super(key: key);

  @override
  _GameEquippedWidgetState createState() => _GameEquippedWidgetState();
}

class _GameEquippedWidgetState extends State<GameEquippedWidget> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameEquippedWidget');
    log.info('Building..');

    return BlocConsumer<DungeonActionCubit, DungeonActionState>(
      listener: (BuildContext context, DungeonActionState state) {
        log.info('listener...');
      },
      builder: (BuildContext context, DungeonActionState state) {
        log.info('Rendering equipped inventory');

        List<Widget> equippedWidgets = [];
        if (state is DungeonActionStatePlaying &&
            state.current.character != null &&
            state.current.character?.equippedObjects != null) {
          var equipped = state.current.character?.equippedObjects;
          for (var i = 0; i < equipped!.length; i++) {
            equippedWidgets.add(
              ObjectButtonWidget(objectName: equipped[i].name),
            );
          }
        }
        return GridView.count(
          crossAxisCount: 4,
          children: equippedWidgets,
        );
      },
    );
  }
}
