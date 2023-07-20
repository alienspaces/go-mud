import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon_character/dungeon_character_cubit.dart';
import 'package:go_mud_client/widgets/game/board/board.dart';
import 'package:go_mud_client/widgets/game/action/panel.dart';
import 'package:go_mud_client/widgets/game/action/narrative.dart';
import 'package:go_mud_client/widgets/game/action/command.dart';
import 'package:go_mud_client/widgets/game/board/location/description/description_container.dart';
import 'package:go_mud_client/widgets/game/card/card.dart';

class GameWidget extends StatefulWidget {
  const GameWidget({Key? key}) : super(key: key);

  @override
  State<GameWidget> createState() => _GameWidgetState();
}

class _GameWidgetState extends State<GameWidget> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameWidget', 'build');
    log.info('Building..');

    return BlocConsumer<DungeonCharacterCubit, DungeonCharacterState>(
      listener: (context, state) {
        log.fine('listener...');
      },
      builder: (BuildContext context, DungeonCharacterState state) {
        log.fine('builder...');

        if (state is DungeonCharacterStateCreate) {
          // ignore: avoid_unnecessary_containers
          return Container(
            child: const Text("GameWidget - Entering"),
          );
        } else if (state is DungeonCharacterStateCreateError) {
          // ignore: avoid_unnecessary_containers
          return Container(
            child: const Text("GameWidget - Error"),
          );
        } else if (state is DungeonCharacterStateCreated) {
          // ignore: avoid_unnecessary_containers
          return Container(
            color: Colors.purple[50],
            // Top level game widget stack
            child: Stack(
              children: <Widget>[
                // ignore: avoid_unnecessary_containers
                Container(
                  child: const Column(
                    children: <Widget>[
                      // Location description
                      Expanded(
                        flex: 3,
                        child: GameLocationDescriptionContainerWidget(),
                      ),
                      // Game board
                      Expanded(
                        flex: 10,
                        child: GameBoardWidget(),
                      ),
                      // Current actions
                      Expanded(
                        flex: 4,
                        child: GameActionPanelWidget(),
                      ),
                      // Current command
                      Expanded(
                        flex: 1,
                        child: GameActionCommandWidget(),
                      ),
                    ],
                  ),
                ),
                // ignore: avoid_unnecessary_containers
                Container(
                  child: const GameCardWidget(),
                ),
                // ignore: avoid_unnecessary_containers
                Container(
                  child: const GameActionNarrativeWidget(),
                )
              ],
            ),
          );
        }

        // ignore: avoid_unnecessary_containers
        return Container(
          child: const Text("Empty"),
        );
      },
    );
  }
}
