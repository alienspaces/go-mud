import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/widgets/game/game_dungeon_grid.dart';

class GameDungeonGridContainerWidget extends StatefulWidget {
  const GameDungeonGridContainerWidget({Key? key}) : super(key: key);

  @override
  _GameDungeonGridContainerWidgetState createState() => _GameDungeonGridContainerWidgetState();
}

class _GameDungeonGridContainerWidgetState extends State<GameDungeonGridContainerWidget> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameDungeonGridWidget');
    log.info('Building..');

    return BlocConsumer<DungeonActionCubit, DungeonActionState>(
      listener: (BuildContext context, DungeonActionState state) {
        log.info('listener...');
      },
      builder: (BuildContext context, DungeonActionState state) {
        if (state is DungeonActionStateCreated) {
          // TODO: Based on the resulting dungeon action command animate the result
          return Stack(
            children: const [
              GameDungeonGridWidget(),
            ],
          );
        }
        return Container();
      },
    );
  }
}
