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
      // Does not render the location grid when there the action
      // is being created or has an error.
      buildWhen: (DungeonActionState prevState, DungeonActionState currState) {
        if (currState is DungeonActionStateError ||
            currState is DungeonActionStateCreating) {
          return false;
        }
        return true;
      },
      builder: (BuildContext context, DungeonActionState state) {
        List<Widget> widgets = [];

        if (state is DungeonActionStateCreated) {
          var actionRec = state.action;

          log.fine(
              'DungeonActionStateCreated - Rendering action ${actionRec.actionCommand}');

          widgets.add(
            GameLocationGridWidget(
              key: UniqueKey(),
              action: actionRec.actionCommand,
              locationData: actionRec.actionLocation,
            ),
          );
        } else if (state is DungeonActionStatePlaying) {
          var currentActionRec = state.currentActionRec;
          var previousActionRec = state.previousActionRec;

          log.fine(
              'DungeonActionStatePlaying - Rendering action ${currentActionRec.actionCommand}');

          if (currentActionRec.actionCommand == 'move' &&
              previousActionRec != null) {
            // Slides the previous location out
            widgets.add(
              GameLocationGridMoveWidget(
                key: UniqueKey(),
                slide: Slide.slideOut,
                direction: state.actionDirection,
                action: currentActionRec.actionCommand,
                locationData: previousActionRec.actionLocation,
              ),
            );
            // Slides the current location in
            widgets.add(
              GameLocationGridMoveWidget(
                key: UniqueKey(),
                slide: Slide.slideIn,
                direction: state.actionDirection,
                action: state.actionCommand,
                locationData: currentActionRec.actionLocation,
              ),
            );
          } else if (currentActionRec.actionCommand == 'look' ||
              currentActionRec.actionCommand == 'attack') {
            widgets.add(
              GameLocationGridWidget(
                key: UniqueKey(),
                action: currentActionRec.actionCommand,
                locationData: currentActionRec.actionLocation,
              ),
            );
            if (currentActionRec.actionTargetLocation != null) {
              log.fine('Rendering look target location');
              widgets.add(
                GameLocationGridLookWidget(
                  key: UniqueKey(),
                  direction: state.actionDirection,
                  action: currentActionRec.actionCommand,
                  locationData: currentActionRec.actionTargetLocation!,
                ),
              );
            }
          } else if (['stash', 'equip', 'drop']
              .contains(currentActionRec.actionCommand)) {
            widgets.add(
              GameLocationGridWidget(
                key: UniqueKey(),
                action: currentActionRec.actionCommand,
                locationData: currentActionRec.actionLocation,
              ),
            );
          }
        }

        log.info('Rendering ${widgets.length} dungeon grid widgets');

        return Stack(
          clipBehavior: Clip.antiAlias,
          children: widgets,
        );
      },
    );
  }
}
