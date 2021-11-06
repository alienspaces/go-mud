import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/widgets/game/game_dungeon_action.dart';
import 'package:go_mud_client/widgets/game/game_dungeon_command.dart';
import 'package:go_mud_client/widgets/game/game_dungeon_description_container.dart';
import 'package:go_mud_client/widgets/game/game_dungeon_grid_container.dart';

class GameDungeonWidget extends StatefulWidget {
  const GameDungeonWidget({Key? key}) : super(key: key);

  @override
  _GameDungeonWidgetState createState() => _GameDungeonWidgetState();
}

class _GameDungeonWidgetState extends State<GameDungeonWidget> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameCharacterWidget');
    log.info('Building..');

    return BlocConsumer<DungeonActionCubit, DungeonActionState>(
      listener: (BuildContext context, DungeonActionState state) {
        log.info('listener...');
      },
      builder: (BuildContext context, DungeonActionState state) {
        if (state is DungeonActionStateCreated) {
          // ignore: avoid_unnecessary_containers
          return Container(
            color: Colors.orange[100],
            child: Column(
              children: const <Widget>[
                GameDungeonDescriptionContainerWidget(),
                GameDungeonGridContainerWidget(),
                GameDungeonActionWidget(),
                GameDungeonCommandWidget(),
              ],
            ),
          );
        }

        // Empty
        return Container();
      },
    );
  }
}
