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
    log.info('Building..');

    return BlocConsumer<DungeonActionCubit, DungeonActionState>(
      listener: (BuildContext context, DungeonActionState state) {
        log.fine('listener...');
      },
      // Do not re-render the location description when there is an error with
      // submitted an action.
      buildWhen: (DungeonActionState prevState, DungeonActionState currState) {
        if (currState is DungeonActionStateError) {
          log.info('Skipping build..');
          return false;
        }
        log.info('Not skipping build..');
        return true;
      },
      builder: (BuildContext context, DungeonActionState state) {
        if (state is DungeonActionStateCreated) {
          log.info('dungeon state is created');
          List<Widget> widgets = [];
          widgets.add(GameLocationDescriptionWidget(
            fade: DescriptionOpacity.fadeIn,
            dungeonActionRecord: state.current,
          ));
          return Stack(
            children: widgets,
          );
        } else if (state is DungeonActionStatePlaying) {
          log.info('dungeon state is playing');
          List<Widget> widgets = [];
          widgets.add(GameLocationDescriptionWidget(
            fade: DescriptionOpacity.fadeOut,
            dungeonActionRecord: state.previous,
          ));
          widgets.add(GameLocationDescriptionWidget(
            fade: DescriptionOpacity.fadeIn,
            dungeonActionRecord: state.current,
          ));
          return Stack(
            children: widgets,
          );
        }
        return Container();
      },
    );
  }
}
