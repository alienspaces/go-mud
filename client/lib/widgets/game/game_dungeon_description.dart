import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/repository/dungeon_action/dungeon_action_repository.dart';

class GameDungeonDescriptionWidget extends StatefulWidget {
  const GameDungeonDescriptionWidget({Key? key}) : super(key: key);

  @override
  _GameDungeonDescriptionWidgetState createState() => _GameDungeonDescriptionWidgetState();
}

typedef DungeonDescriptionMemberFunction = Widget Function(DungeonActionRecord record, String key);

class _GameDungeonDescriptionWidgetState extends State<GameDungeonDescriptionWidget> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameDungeonDescriptionWidget');
    log.info('Building..');

    return BlocConsumer<DungeonActionCubit, DungeonActionState>(
      listener: (BuildContext context, DungeonActionState state) {
        log.info('listener...');
      },
      builder: (BuildContext context, DungeonActionState state) {
        if (state is DungeonActionStateCreated) {
          // ignore: avoid_unnecessary_containers
          return Container(
            child: Column(
              children: [
                Container(
                  margin: const EdgeInsets.fromLTRB(5, 10, 5, 5),
                  child: Text('${state.dungeonActionRecord?.location.name}',
                      style: Theme.of(context).textTheme.headline5),
                ),
                Container(
                  margin: const EdgeInsets.fromLTRB(5, 5, 5, 10),
                  child: Text('${state.dungeonActionRecord?.location.description}'),
                ),
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
