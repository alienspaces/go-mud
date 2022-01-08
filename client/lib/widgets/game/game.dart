import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon/dungeon_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_command/dungeon_command_cubit.dart';
import 'package:go_mud_client/cubit/character/character_cubit.dart';
import 'package:go_mud_client/widgets/game/board/board.dart';
import 'package:go_mud_client/widgets/game/action/action.dart';
import 'package:go_mud_client/widgets/game/command/command.dart';
import 'package:go_mud_client/widgets/game/board/location/description/description_container.dart';
// import 'package:go_mud_client/widgets/game/location/location.dart';

class GameWidget extends StatefulWidget {
  const GameWidget({Key? key}) : super(key: key);

  @override
  _GameWidgetState createState() => _GameWidgetState();
}

class _GameWidgetState extends State<GameWidget> {
  @override
  void initState() {
    final log = getLogger('HomeContainerWidget');
    log.info('Initialising state..');

    super.initState();

    _initAction(context);
  }

  void _initAction(BuildContext context) {
    final log = getLogger('GameWidget');
    log.info('Initialising action..');

    final dungeonCubit = BlocProvider.of<DungeonCubit>(context);
    if (dungeonCubit.dungeonRecord == null) {
      log.warning(
          'Dungeon cubit missing dungeon record, cannot initialise action');
      return;
    }

    final characterCubit = BlocProvider.of<CharacterCubit>(context);
    if (characterCubit.characterRecord == null) {
      log.warning(
          'Character cubit missing character record, cannot initialise action');
      return;
    }

    final dungeonCommandCubit = BlocProvider.of<DungeonCommandCubit>(context);
    dungeonCommandCubit.selectAction(
      'look',
    );

    final dungeonActionCubit = BlocProvider.of<DungeonActionCubit>(context);
    dungeonActionCubit.createAction(
      dungeonCubit.dungeonRecord!.id,
      characterCubit.characterRecord!.id,
      dungeonCommandCubit.command(),
    );

    dungeonCommandCubit.unselectAction();
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('Game');
    log.info('Building..');

    // ignore: avoid_unnecessary_containers
    return Container(
      color: Colors.yellow[100],
      child: Column(
        children: const <Widget>[
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
            child: GameActionWidget(),
          ),
          // Current command
          Expanded(
            flex: 1,
            child: GameCommandWidget(),
          ),
        ],
      ),
    );
  }
}
