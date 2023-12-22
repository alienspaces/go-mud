import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/navigation.dart';
// import 'package:go_mud_client/cubit/character/character_cubit.dart';
import 'package:go_mud_client/cubit/character/character_cubit.dart';
import 'package:go_mud_client/widgets/game/board/board.dart';
import 'package:go_mud_client/widgets/game/action/panel.dart';
import 'package:go_mud_client/widgets/game/action/narrative.dart';
import 'package:go_mud_client/widgets/game/action/command.dart';
import 'package:go_mud_client/widgets/game/board/location/description/description_container.dart';
import 'package:go_mud_client/widgets/game/card/card.dart';

class GameContainerWidget extends StatefulWidget {
  final NavigationCallbacks callbacks;

  const GameContainerWidget({
    Key? key,
    required this.callbacks,
  }) : super(key: key);

  @override
  State<GameContainerWidget> createState() => _GameContainerWidgetState();
}

class _GameContainerWidgetState extends State<GameContainerWidget> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameContainerWidget', 'build');
    log.fine('Building..');

    // final characterCubit = BlocProvider.of<CharacterCubit>(
    //   context,
    //   listen: true,
    // );

    // var characterRecord = characterCubit.characterRecord;
    // if (characterRecord == null) {
    //   log.warning("character record is null, cannot display game");
    //   return const SizedBox.shrink();
    // }

    // if (characterRecord.dungeonID == null) {
    //   log.warning("character record dungeon id is null, cannot display game");
    //   return const SizedBox.shrink();
    // }

    return BlocConsumer<CharacterCubit, CharacterState>(
      listener: (context, state) {
        log.fine('listener...');
      },
      builder: (BuildContext context, CharacterState state) {
        log.fine('builder...');

        if (state is CharacterStateEntering) {
          // ignore: avoid_unnecessary_containers
          return Container(
            color: Colors.purple[50],
            child: const Text("GameContainerWidget - Entering"),
          );
        } else if (state is CharacterStateError) {
          // ignore: avoid_unnecessary_containers
          return Container(
            color: Colors.purple[50],
            child: const Text("GameContainerWidget - Error"),
          );
        } else if (state is CharacterStateEntered) {
          log.warning('Dungeon character cubit emitted entered state');

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

        log.warning('Dungeon character cubit emitted initial state');

        return const SizedBox.shrink();
      },
    );
  }
}
