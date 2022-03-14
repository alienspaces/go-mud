import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';

import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';

import 'package:go_mud_client/widgets/game/button/object_button.dart';

class GameStashedWidget extends StatefulWidget {
  const GameStashedWidget({Key? key}) : super(key: key);

  @override
  _GameStashedWidgetState createState() => _GameStashedWidgetState();
}

class _GameStashedWidgetState extends State<GameStashedWidget> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameStashedWidget');
    log.fine('Building..');

    return BlocConsumer<DungeonActionCubit, DungeonActionState>(
      listener: (BuildContext context, DungeonActionState state) {
        log.fine('listener...');
      },
      builder: (BuildContext context, DungeonActionState state) {
        log.fine('Rendering stashed inventory');

        List<Widget> stashedWidgets = [];
        if (state is DungeonActionStatePlaying &&
            state.current.character != null &&
            state.current.character?.stashedObjects != null) {
          var stashed = state.current.character?.stashedObjects;
          for (var i = 0; i < stashed!.length; i++) {
            stashedWidgets.add(
              ObjectButtonWidget(objectName: stashed[i].name),
            );
          }
        }
        return GridView.count(
          crossAxisCount: 5,
          children: stashedWidgets,
        );
      },
    );
  }
}
