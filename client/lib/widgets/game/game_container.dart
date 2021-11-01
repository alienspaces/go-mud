import 'package:go_mud_client/widgets/game/game_character.dart';
import 'package:go_mud_client/widgets/game/game_dungeon.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon/dungeon_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_command/dungeon_command_cubit.dart';
import 'package:go_mud_client/cubit/character/character_cubit.dart';

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
      log.warning('Dungeon cubit missing dungeon record, cannot initialise action');
      return;
    }

    final characterCubit = BlocProvider.of<CharacterCubit>(context);
    if (characterCubit.characterRecord == null) {
      log.warning('Character cubit missing character record, cannot initialise action');
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
      child: Column(
        children: const <Widget>[
          GameCharacterWidget(),
          GameDungeonWidget(),
        ],
      ),
    );
  }
}
