import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon/dungeon_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_command/dungeon_command_cubit.dart';
import 'package:go_mud_client/cubit/character/character_cubit.dart';
import 'package:go_mud_client/widgets/game/game_character.dart';

// import 'package:go_mud_client/widgets/game/game_dungeon.dart';
import 'package:go_mud_client/widgets/game/game_dungeon_action.dart';
import 'package:go_mud_client/widgets/game/game_dungeon_command.dart';
import 'package:go_mud_client/widgets/game/game_dungeon_description_container.dart';
import 'package:go_mud_client/widgets/game/game_dungeon_container.dart';

class GameContainerWidget extends StatefulWidget {
  const GameContainerWidget({Key? key}) : super(key: key);

  @override
  _GameContainerWidgetState createState() => _GameContainerWidgetState();
}

class _GameContainerWidgetState extends State<GameContainerWidget> {
  @override
  void initState() {
    final log = getLogger('HomeContainerWidget');
    log.info('Initialising state..');

    super.initState();

    _initAction(context);
  }

  void _initAction(BuildContext context) {
    final log = getLogger('GameContainerWidget');
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
    final log = getLogger('GameContainer');
    log.info('Building..');

    // ignore: avoid_unnecessary_containers
    return Container(
      color: Colors.yellow[100],
      child: Column(
        children: <Widget>[
          // Character
          const Expanded(
            flex: 5,
            child: GameCharacterWidget(),
          ),
          // Location description container
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
