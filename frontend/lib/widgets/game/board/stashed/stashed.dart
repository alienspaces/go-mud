import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/widgets/game/button/object_button.dart';

class GameStashedWidget extends StatefulWidget {
  const GameStashedWidget({Key? key}) : super(key: key);

  @override
  State<GameStashedWidget> createState() => _GameStashedWidgetState();
}

class _GameStashedWidgetState extends State<GameStashedWidget> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameStashedWidget', 'build');
    log.fine('Building..');

    return BlocConsumer<DungeonActionCubit, DungeonActionState>(
      listener: (BuildContext context, DungeonActionState state) {
        log.fine('listener...');
      },
      builder: (BuildContext context, DungeonActionState state) {
        log.fine('Rendering stashed inventory');

        List<Widget> stashedWidgets = [];
        if (state is DungeonActionStatePlaying) {
          var actionCharacter = state.currentActionRec.actionCharacter;
          if (actionCharacter != null &&
              actionCharacter.characterStashedObjects != null) {
            var stashed = actionCharacter.characterStashedObjects;
            for (var i = 0; i < stashed!.length; i++) {
              stashedWidgets.add(
                ObjectButtonWidget(name: stashed[i].objectName),
              );
            }
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
