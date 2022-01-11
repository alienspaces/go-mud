import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/widgets/game/board/location/grid/grid.dart';
import 'package:go_mud_client/widgets/game/board/location/grid/grid_move.dart';
import 'package:go_mud_client/widgets/game/board/location/grid/grid_look.dart';

class BoardLocationWidget extends StatelessWidget {
  const BoardLocationWidget({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameLocationWidget');
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
              GameLocationGridWidget(
                key: UniqueKey(),
                action: null,
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
            GameLocationGridWidget(
              key: UniqueKey(),
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
              GameLocationGridMoveWidget(
                key: UniqueKey(),
                slide: Slide.slideOut,
                direction: state.direction,
                action: state.action,
                locationData: state.previous.location,
              ),
            );
            widgets.add(
              GameLocationGridMoveWidget(
                key: UniqueKey(),
                slide: Slide.slideIn,
                direction: state.direction,
                action: state.action,
                locationData: state.current.location,
              ),
            );
          } else if (dungeonActionRecord.command == 'look') {
            widgets.add(
              GameLocationGridWidget(
                key: UniqueKey(),
                action: state.action,
                locationData: state.current.location,
              ),
            );
            if (dungeonActionRecord.targetLocation != null) {
              log.info('Rendering look target location');
              widgets.add(
                GameLocationGridLookWidget(
                  key: UniqueKey(),
                  direction: state.direction,
                  action: state.action,
                  locationData: state.current.targetLocation!,
                ),
              );
            }
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
