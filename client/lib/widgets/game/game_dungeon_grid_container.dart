import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/widgets/game/game_dungeon_sliding_grid.dart';

class GameDungeonGridContainerWidget extends StatefulWidget {
  const GameDungeonGridContainerWidget({Key? key}) : super(key: key);

  @override
  _GameDungeonGridContainerWidgetState createState() => _GameDungeonGridContainerWidgetState();
}

class _GameDungeonGridContainerWidgetState extends State<GameDungeonGridContainerWidget> {
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
        if (state is DungeonActionStateCreated) {
          var dungeonActionRecord = state.current;

          log.info(
              'DungeonActionStateCreated - Rendering action ${dungeonActionRecord.action.command}');

          if (dungeonActionRecord.action.command == 'move') {
            log.info('DungeonActionStateCreated - Rendering move');
            widgets.add(
              GameDungeonSlidingGridWidget(
                key: UniqueKey(),
                slide: Slide.slideNone,
                action: state.action,
                dungeonActionRecord: dungeonActionRecord,
              ),
            );
          } else if (dungeonActionRecord.action.command == 'look') {
            log.info('DungeonActionStateCreated - Rendering look');
            widgets.add(
              GameDungeonSlidingGridWidget(
                key: UniqueKey(),
                slide: Slide.slideNone,
                action: state.action,
                dungeonActionRecord: dungeonActionRecord,
              ),
            );
          }
        } else if (state is DungeonActionStateCreating) {
          var dungeonActionRecord = state.current;

          log.info(
              'DungeonActionStateCreating - Rendering command ${dungeonActionRecord?.action.command}');

          if (dungeonActionRecord != null) {
            log.info('DungeonActionStateCreating - Rendering move');
            widgets.add(
              GameDungeonSlidingGridWidget(
                key: UniqueKey(),
                slide: Slide.slideNone,
                dungeonActionRecord: dungeonActionRecord,
              ),
            );
          } else {
            log.info('DungeonActionStateCreating - Record is null');
            widgets.add(
              Container(
                color: Colors.blueAccent,
                child: const Text('Loading'),
              ),
            );
          }
        } else if (state is DungeonActionStatePlaying) {
          var dungeonActionRecord = state.current;

          log.info(
              'DungeonActionStatePlaying - Rendering action ${dungeonActionRecord.action.command}');

          if (dungeonActionRecord.action.command == 'move') {
            log.info('DungeonActionStatePlaying - Rendering move');
            widgets.add(
              GameDungeonSlidingGridWidget(
                key: UniqueKey(),
                slide: Slide.slideOut,
                direction: state.direction,
                action: state.action,
                dungeonActionRecord: state.previous,
              ),
            );
            widgets.add(
              GameDungeonSlidingGridWidget(
                key: UniqueKey(),
                slide: Slide.slideIn,
                direction: state.direction,
                action: state.action,
                dungeonActionRecord: state.current,
              ),
            );
          } else if (dungeonActionRecord.action.command == 'look') {
            log.info('DungeonActionStatePlaying - Rendering look');
            widgets.add(
              GameDungeonSlidingGridWidget(
                key: UniqueKey(),
                slide: Slide.slideNone,
                action: state.action,
                dungeonActionRecord: dungeonActionRecord,
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
