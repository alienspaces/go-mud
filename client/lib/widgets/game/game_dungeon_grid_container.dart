import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/widgets/game/game_dungeon_grid.dart';

class GameDungeonGridContainerWidget extends StatefulWidget {
  const GameDungeonGridContainerWidget({Key? key}) : super(key: key);

  @override
  _GameDungeonGridContainerWidgetState createState() => _GameDungeonGridContainerWidgetState();
}

class _GameDungeonGridContainerWidgetState extends State<GameDungeonGridContainerWidget> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameDungeonGridWidget');
    log.info('Building..');

    return BlocConsumer<DungeonActionCubit, DungeonActionState>(
      listener: (BuildContext context, DungeonActionState state) {
        log.info('listener...');
      },
      builder: (BuildContext context, DungeonActionState state) {
        if (state is DungeonActionStateCreated) {
          List<Widget> widgets = [];
          // TODO: "play" actions here?
          var dungeonActionRecord = state.dungeonActionRecord;

          if (dungeonActionRecord != null) {
            log.info('Animating action command ${dungeonActionRecord.action.command}');
            if (dungeonActionRecord.action.command == 'move') {
              log.info('Animating move action');
              widgets.add(
                GameDungeonGridWidget(
                  scroll: DungeonGridScroll.scrollNone,
                  dungeonActionRecord: dungeonActionRecord,
                ),
              );
            } else if (dungeonActionRecord.action.command == 'look') {
              log.info('Animating look action');
              widgets.add(
                GameDungeonGridWidget(
                  scroll: DungeonGridScroll.scrollNone,
                  dungeonActionRecord: dungeonActionRecord,
                ),
              );
            }
          }
          return Stack(
            children: widgets,
          );
        } else if (state is DungeonActionStatePlaying) {
          List<Widget> widgets = [];
          // TODO: "play" actions here?
          var dungeonActionRecord = state.current;

          log.info('Animating action command ${dungeonActionRecord.action.command}');
          if (dungeonActionRecord.action.command == 'move') {
            log.info('Animating move action');
            widgets.add(
              GameDungeonGridWidget(
                scroll: DungeonGridScroll.scrollOut,
                dungeonActionRecord: state.previous,
              ),
            );
            widgets.add(
              GameDungeonGridWidget(
                scroll: DungeonGridScroll.scrollIn,
                dungeonActionRecord: state.current,
              ),
            );
          } else if (dungeonActionRecord.action.command == 'look') {
            log.info('Animating look action');
            widgets.add(
              GameDungeonGridWidget(
                scroll: DungeonGridScroll.scrollNone,
                dungeonActionRecord: dungeonActionRecord,
              ),
            );
          }
          return Stack(
            clipBehavior: Clip.hardEdge,
            children: widgets,
          );
        }

        return Container();
      },
    );
  }
}
