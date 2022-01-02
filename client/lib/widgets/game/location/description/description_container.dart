import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/widgets/game/location/description/description.dart';

class GameLocationDescriptionContainerWidget extends StatefulWidget {
  const GameLocationDescriptionContainerWidget({Key? key}) : super(key: key);

  @override
  _GameLocationDescriptionContainerWidgetState createState() =>
      _GameLocationDescriptionContainerWidgetState();
}

class _GameLocationDescriptionContainerWidgetState
    extends State<GameLocationDescriptionContainerWidget> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameLocationDescriptionContainerWidget');
    log.info('Building..');

    return BlocConsumer<DungeonActionCubit, DungeonActionState>(
      listener: (BuildContext context, DungeonActionState state) {
        log.info('listener...');
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
