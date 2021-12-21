import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/widgets/game/game_dungeon_move_grid.dart';
import 'package:go_mud_client/widgets/game/game_dungeon_look_grid.dart';

class GameDungeonGridContainerWidget extends StatefulWidget {
  const GameDungeonGridContainerWidget({Key? key}) : super(key: key);

  @override
  _GameDungeonGridContainerWidgetState createState() =>
      _GameDungeonGridContainerWidgetState();
}

class _GameDungeonGridContainerWidgetState
    extends State<GameDungeonGridContainerWidget> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameDungeonGridContainerWidget');
    log.info('Building..');

    return BlocConsumer<DungeonActionCubit, DungeonActionState>(
      listener: (BuildContext context, DungeonActionState state) {
        log.info('listener...');
      },
      builder: (BuildContext context, DungeonActionState state) {
        List<Widget> widgets = [];

        // Creating state is emitted with every action
        if (state is DungeonActionStateCreating) {
          var dungeonActionRecord = state.current;

          if (dungeonActionRecord != null) {
            log.info(
                'DungeonActionStateCreating - Rendering command ${dungeonActionRecord.command}');
            widgets.add(
              GameDungeonMoveGridWidget(
                key: UniqueKey(),
                slide: Slide.slideNone,
                locationData: dungeonActionRecord.location,
              ),
            );
          } else {
            log.info(
                'DungeonActionStateCreating - Rendering loading container..');
            widgets.add(
              Container(
                color: Colors.blueAccent,
                child: const Text('Loading'),
              ),
            );
          }
        }
        // Created state is emitted with every action
        else if (state is DungeonActionStateCreated) {
          var dungeonActionRecord = state.current;

          log.info(
              'DungeonActionStateCreated - Rendering action ${dungeonActionRecord.command}');

          widgets.add(
            GameDungeonMoveGridWidget(
              key: UniqueKey(),
              slide: Slide.slideNone,
              action: state.action,
              locationData: dungeonActionRecord.location,
            ),
          );
        }
        // Playing state is emitted only when there is at least one previous action
        else if (state is DungeonActionStatePlaying) {
          var dungeonActionRecord = state.current;

          log.info(
              'DungeonActionStatePlaying - Rendering action ${dungeonActionRecord.command}');

          if (dungeonActionRecord.command == 'move') {
            widgets.add(
              GameDungeonMoveGridWidget(
                key: UniqueKey(),
                slide: Slide.slideOut,
                direction: state.direction,
                action: state.action,
                locationData: state.previous.location,
              ),
            );
            widgets.add(
              GameDungeonMoveGridWidget(
                key: UniqueKey(),
                slide: Slide.slideIn,
                direction: state.direction,
                action: state.action,
                locationData: state.current.location,
              ),
            );
          } else if (dungeonActionRecord.command == 'look') {
            widgets.add(
              GameDungeonMoveGridWidget(
                key: UniqueKey(),
                slide: Slide.slideNone,
                direction: state.direction,
                action: state.action,
                locationData: state.current.location,
              ),
            );
            if (state.current.targetLocation != null) {
              log.info('Rendering look target location');
              widgets.add(
                GameDungeonLookGridWidget(
                  key: UniqueKey(),
                  direction: state.direction,
                  action: state.action,
                  locationData: state.current.targetLocation!,
                ),
              );
            } else if (state.current.targetCharacter != null) {
              widgets.add(
                Container(
                  padding: const EdgeInsets.all(5),
                  child: const Text("Looking character"),
                ),
              );
            }
          } else if (state.current.targetMonster != null) {
            widgets.add(
              Container(
                padding: const EdgeInsets.all(5),
                child: const Text("Looking monster"),
              ),
            );
          } else if (state.current.targetObject != null) {
            widgets.add(
              Container(
                padding: const EdgeInsets.all(5),
                child: const Text("Looking object"),
              ),
            );
          }
        }

        log.info('Rendering ${widgets.length} dungeon grid panels');

        return Stack(
          clipBehavior: Clip.antiAlias,
          children: widgets,
        );
      },
    );
  }
}
