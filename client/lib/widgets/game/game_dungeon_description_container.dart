import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/widgets/game/game_dungeon_description.dart';

class GameDungeonDescriptionContainerWidget extends StatefulWidget {
  const GameDungeonDescriptionContainerWidget({Key? key}) : super(key: key);

  @override
  _GameDungeonDescriptionContainerWidgetState createState() =>
      _GameDungeonDescriptionContainerWidgetState();
}

class _GameDungeonDescriptionContainerWidgetState
    extends State<GameDungeonDescriptionContainerWidget> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameDungeonDescriptionContainerWidget');
    log.info('Building..');

    return BlocConsumer<DungeonActionCubit, DungeonActionState>(
      listener: (BuildContext context, DungeonActionState state) {
        log.info('listener...');
      },
      builder: (BuildContext context, DungeonActionState state) {
        if (state is DungeonActionStateCreated) {
          log.info('dungeon state is created');
          List<Widget> widgets = [];
          widgets.add(GameDungeonDescriptionWidget(
            fade: DescriptionOpacity.fadeIn,
            dungeonActionRecord: state.current,
          ));
          return Stack(
            children: widgets,
          );
        } else if (state is DungeonActionStatePlaying) {
          log.info('dungeon state is playing');
          List<Widget> widgets = [];
          widgets.add(GameDungeonDescriptionWidget(
            fade: DescriptionOpacity.fadeOut,
            dungeonActionRecord: state.previous,
          ));
          widgets.add(GameDungeonDescriptionWidget(
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
