import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/widgets/game/game_dungeon_action.dart';
import 'package:go_mud_client/widgets/game/game_dungeon_command.dart';
import 'package:go_mud_client/widgets/game/game_dungeon_description_container.dart';
import 'package:go_mud_client/widgets/game/game_dungeon_container.dart';

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

    return Container(
      color: Colors.orange[100],
      child: Column(
        children: <Widget>[
          // Location description
          const Expanded(
            flex: 3,
            child: GameDungeonDescriptionContainerWidget(),
          ),
          // Location container
          Expanded(
            flex: 10,
            child: Container(
              decoration: BoxDecoration(color: Colors.orange[100]),
              clipBehavior: Clip.antiAlias,
              child: const GameDungeonContainerWidget(),
            ),
          ),
          // Location actions
          const Expanded(
            flex: 4,
            child: GameDungeonActionWidget(),
          ),
          // Current command
          const Expanded(
            flex: 1,
            child: GameDungeonCommandWidget(),
          ),
        ],
      ),
    );
  }
}
