import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/widgets/game/board/location/description/description.dart';

class GameLocationDescriptionContainerWidget extends StatelessWidget {
  const GameLocationDescriptionContainerWidget({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameLocationDescriptionContainerWidget', 'build');
    log.fine('Building..');

    return BlocConsumer<DungeonActionCubit, DungeonActionState>(
      listener: (BuildContext context, DungeonActionState state) {
        log.fine('listener...');
      },
      // Does not render the location grid when there the action
      // is being created or has an error.
      buildWhen: (DungeonActionState prevState, DungeonActionState currState) {
        if (currState is DungeonActionStateError ||
            currState is DungeonActionStateCreating ||
            currState is DungeonActionStatePlayingOther) {
          return false;
        }

        if (currState is DungeonActionStatePlaying &&
            currState.currentActionRec.actionLocation.locationName !=
                currState.currentActionRec.actionTargetLocation!.locationName) {
          return false;
        }

        return true;
      },
      builder: (BuildContext context, DungeonActionState state) {
        List<Widget> widgets = [];

        if (state is DungeonActionStateCreated) {
          log.warning('dungeon state is created');
          widgets.add(GameLocationDescriptionWidget(
            fade: DescriptionOpacity.fadeIn,
            dungeonActionRecord: state.action,
          ));
        } else if (state is DungeonActionStatePlaying) {
          log.warning('dungeon state is playing');
          if (state.previousActionRec != null) {
            widgets.add(GameLocationDescriptionWidget(
              fade: DescriptionOpacity.fadeOut,
              dungeonActionRecord: state.previousActionRec!,
            ));
          }
          widgets.add(GameLocationDescriptionWidget(
            fade: DescriptionOpacity.fadeIn,
            dungeonActionRecord: state.currentActionRec,
          ));
        }

        log.fine('Rendering ${widgets.length} dungeon description widgets');

        return Stack(
          children: widgets,
        );
      },
    );
  }
}
