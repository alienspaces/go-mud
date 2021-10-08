import 'package:go_mud_client/repository/dungeon_action/dungeon_action_repository.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon/dungeon_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/cubit/character/character_cubit.dart';

class GameDungeonGridWidget extends StatefulWidget {
  const GameDungeonGridWidget({Key? key}) : super(key: key);

  @override
  _GameDungeonGridWidgetState createState() => _GameDungeonGridWidgetState();
}

typedef DungeonGridMemberFunction = Widget Function(DungeonActionRecord record, String key);

class _GameDungeonGridWidgetState extends State<GameDungeonGridWidget> {
  double gridMemberWidth = 50;
  double gridMemberHeight = 50;

  Map<String, String> directionLabelMap = {
    'north': 'N',
    'northeast': 'NE',
    'east': 'E',
    'southeast': 'SE',
    'south': 'S',
    'southwest': 'SW',
    'west': 'W',
    'northwest': 'NW',
    'up': 'U',
    'down': 'D',
  };

  // Direction widget
  Widget directionWidget(BuildContext context, DungeonActionRecord record, String direction) {
    if (record.location.directions.contains(direction)) {
      return Container(
        margin: const EdgeInsets.all(2),
        child: ElevatedButton(
          onPressed: () {
            final log = getLogger('GameDungeonGridWidget');

            final dungeonCubit = BlocProvider.of<DungeonCubit>(context);
            if (dungeonCubit.dungeonRecord == null) {
              log.warning(
                'onPressed - Dungeon cubit missing dungeon record, cannot initialise action',
              );
              return;
            }
            // TODO: Select action prior to selecting the direction.
            _submitAction(context, 'move $direction');
          },
          child: Text('${directionLabelMap[direction]}'),
        ),
      );
    }
    return emptyWidget('${directionLabelMap[direction]}');
  }

  // Object widget
  Widget objectWidget(BuildContext context, DungeonActionRecord record, int idx) {
    return emptyWidget('O$idx');
  }

  // Character widget
  Widget characterWidget(BuildContext context, DungeonActionRecord record, int idx) {
    return emptyWidget('C$idx');
  }

  // Monster widget
  Widget monsterWidget(BuildContext context, DungeonActionRecord record, int idx) {
    return emptyWidget('M$idx');
  }

  // Empty widget
  Widget emptyWidget(String label) {
    return Container(
      width: gridMemberWidth,
      height: gridMemberHeight,
      alignment: Alignment.center,
      margin: const EdgeInsets.all(2),
      decoration: BoxDecoration(
        color: const Color(0xFFD4D4D4),
        border: Border.all(
          color: const Color(0xFFD4D4D4),
        ),
        borderRadius: const BorderRadius.all(Radius.circular(5)),
      ),
      child: Text(label),
    );
  }

  List<Widget> generateGrid(BuildContext context) {
    final log = getLogger('CharacterCreateWidget');

    final dungeonActionCubit = BlocProvider.of<DungeonActionCubit>(context);
    if (dungeonActionCubit.dungeonActionRecord == null) {
      log.warning('Dungeon cubit dungeon record is null, cannot create character');
      return [];
    }

    // TODO: Objects, monsters and characters should be randomly scattered through the
    // room but not change position in grid with a widget rebuild..

    List<Widget Function()> dunegonGridMemberFunctions = [
      // Top Row
      () => directionWidget(context, dungeonActionCubit.dungeonActionRecord!, 'northwest'),
      () => objectWidget(context, dungeonActionCubit.dungeonActionRecord!, 0),
      () => directionWidget(context, dungeonActionCubit.dungeonActionRecord!, 'north'),
      () => Container(),
      () => directionWidget(context, dungeonActionCubit.dungeonActionRecord!, 'northeast'),
      // Second Row
      () => objectWidget(context, dungeonActionCubit.dungeonActionRecord!, 1),
      () => objectWidget(context, dungeonActionCubit.dungeonActionRecord!, 2),
      () => directionWidget(context, dungeonActionCubit.dungeonActionRecord!, 'up'),
      () => Container(),
      () => Container(),
      // Third Row
      () => directionWidget(context, dungeonActionCubit.dungeonActionRecord!, 'west'),
      () => objectWidget(context, dungeonActionCubit.dungeonActionRecord!, 3),
      () => Container(),
      () => Container(),
      () => directionWidget(context, dungeonActionCubit.dungeonActionRecord!, 'east'),
      // Fourth Row
      () => objectWidget(context, dungeonActionCubit.dungeonActionRecord!, 4),
      () => objectWidget(context, dungeonActionCubit.dungeonActionRecord!, 5),
      () => directionWidget(context, dungeonActionCubit.dungeonActionRecord!, 'down'),
      () => Container(),
      () => Container(),
      // Bottom Row
      () => directionWidget(context, dungeonActionCubit.dungeonActionRecord!, 'southwest'),
      () => Container(),
      () => directionWidget(context, dungeonActionCubit.dungeonActionRecord!, 'south'),
      () => Container(),
      () => directionWidget(context, dungeonActionCubit.dungeonActionRecord!, 'southeast'),
    ];

    List<Widget> gridWidgets = [];
    for (var gridMemberFunction in dunegonGridMemberFunctions) {
      gridWidgets.add(gridMemberFunction());
    }

    return gridWidgets;
  }

  void _submitAction(BuildContext context, String action) {
    final log = getLogger('GameDungeonGridWidget');
    log.info('Submitting move action..');

    final dungeonCubit = BlocProvider.of<DungeonCubit>(context);
    if (dungeonCubit.dungeonRecord == null) {
      log.warning('Dungeon cubit missing dungeon record, cannot initialise action');
      return;
    }

    final characterCubit = BlocProvider.of<CharacterCubit>(context);
    if (characterCubit.characterRecord == null) {
      log.warning('Character cubit missing character record, cannot initialise action');
      return;
    }

    final dungeonActionCubit = BlocProvider.of<DungeonActionCubit>(context);
    dungeonActionCubit.createAction(
      dungeonCubit.dungeonRecord!.id,
      characterCubit.characterRecord!.id,
      action,
    );
  }

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
          return Container(
            decoration: BoxDecoration(
              color: const Color(0xFFDEDEDE),
              border: Border.all(
                color: const Color(0xFFDEDEDE),
              ),
              borderRadius: const BorderRadius.all(Radius.circular(5)),
            ),
            padding: const EdgeInsets.all(1),
            margin: const EdgeInsets.all(5),
            width: gridMemberWidth * 5,
            height: gridMemberHeight * 5,
            child: GridView.count(
              crossAxisCount: 5,
              children: generateGrid(context),
            ),
          );
        }

        // Empty
        return Container();
      },
    );
  }
}
