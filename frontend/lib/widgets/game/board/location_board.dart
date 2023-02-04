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
    final log = getLogger('BoardLocationWidget', 'build');
    log.fine('Building..');

    return BlocConsumer<DungeonActionCubit, DungeonActionState>(
      listener: (BuildContext context, DungeonActionState state) {
        log.fine('listener...');
      },
      // Do not re-render the location grid when there is an error with
      // submitted an action.
      buildWhen: (DungeonActionState prevState, DungeonActionState currState) {
        if (currState is DungeonActionStateError) {
          return false;
        }
        return true;
      },
      builder: (BuildContext context, DungeonActionState state) {
        List<Widget> widgets = [];

        // Creating state is emitted with every action
        if (state is DungeonActionStateCreating) {
          var dungeonActionRecord = state.current;

          if (dungeonActionRecord != null) {
            log.fine(
                'DungeonActionStateCreating - Rendering command ${dungeonActionRecord.actionCommand}');
            widgets.add(
              GameLocationGridWidget(
                key: UniqueKey(),
                action: null,
                locationData: dungeonActionRecord.actionLocation,
              ),
            );
          } else {
            log.fine(
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

          log.fine(
              'DungeonActionStateCreated - Rendering action ${dungeonActionRecord.actionCommand}');

          widgets.add(
            GameLocationGridWidget(
              key: UniqueKey(),
              action: state.action,
              locationData: dungeonActionRecord.actionLocation,
            ),
          );
        }
        // Playing state is emitted only when there is at least one previous action
        else if (state is DungeonActionStatePlaying) {
          var dungeonActionRecord = state.current;

          log.fine(
              'DungeonActionStatePlaying - Rendering action ${dungeonActionRecord.actionCommand}');

          if (dungeonActionRecord.actionCommand == 'move') {
            widgets.add(
              GameLocationGridMoveWidget(
                key: UniqueKey(),
                slide: Slide.slideOut,
                direction: state.direction,
                action: state.action,
                locationData: state.previous.actionLocation,
              ),
            );
            widgets.add(
              GameLocationGridMoveWidget(
                key: UniqueKey(),
                slide: Slide.slideIn,
                direction: state.direction,
                action: state.action,
                locationData: state.current.actionLocation,
              ),
            );
          } else if (dungeonActionRecord.actionCommand == 'look' ||
              dungeonActionRecord.actionCommand == 'attack') {
            widgets.add(
              GameLocationGridWidget(
                key: UniqueKey(),
                action: state.action,
                locationData: state.current.actionLocation,
              ),
            );
            if (dungeonActionRecord.actionTargetLocation != null) {
              log.fine('Rendering look target location');
              widgets.add(
                GameLocationGridLookWidget(
                  key: UniqueKey(),
                  direction: state.direction,
                  action: state.action,
                  locationData: state.current.actionTargetLocation!,
                ),
              );
            }
          } else if (['stash', 'equip', 'drop']
              .contains(dungeonActionRecord.actionCommand)) {
            widgets.add(
              GameLocationGridWidget(
                key: UniqueKey(),
                action: state.action,
                locationData: state.current.actionLocation,
              ),
            );
          }
        }

        log.fine('Rendering ${widgets.length} dungeon grid panels');

        return Stack(
          clipBehavior: Clip.antiAlias,
          children: widgets,
        );
      },
    );
  }
}
