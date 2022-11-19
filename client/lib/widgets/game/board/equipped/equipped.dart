import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';

import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';

import 'package:go_mud_client/widgets/game/button/object_button.dart';

class GameEquippedWidget extends StatefulWidget {
  const GameEquippedWidget({Key? key}) : super(key: key);

  @override
  State<GameEquippedWidget> createState() => _GameEquippedWidgetState();
}

class _GameEquippedWidgetState extends State<GameEquippedWidget> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameEquippedWidget');
    log.fine('Building..');

    return BlocConsumer<DungeonActionCubit, DungeonActionState>(
      listener: (BuildContext context, DungeonActionState state) {
        log.fine('listener...');
      },
      builder: (BuildContext context, DungeonActionState state) {
        log.fine('Rendering equipped inventory');

        List<Widget> equippedWidgets = [];
        if (state is DungeonActionStatePlaying &&
            state.current.actionCharacter != null &&
            state.current.actionCharacter?.characterEquippedObjects != null) {
          var equipped =
              state.current.actionCharacter?.characterEquippedObjects;
          for (var i = 0; i < equipped!.length; i++) {
            equippedWidgets.add(
              ObjectButtonWidget(objectName: equipped[i].objectName),
            );
          }
        }
        return GridView.count(
          crossAxisCount: 5,
          children: equippedWidgets,
        );
      },
    );
  }
}
