import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/location.dart';
import 'package:go_mud_client/style.dart';
import 'package:go_mud_client/cubit/dungeon/dungeon_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_command/dungeon_command_cubit.dart';
import 'package:go_mud_client/cubit/character/character_cubit.dart';
import 'package:go_mud_client/repository/dungeon_action/dungeon_action_repository.dart';

class GameLocationGridWidget extends StatefulWidget {
  final LocationData locationData;
  final String? action;
  final String? direction;
  final bool readonly;

  const GameLocationGridWidget({
    Key? key,
    required this.locationData,
    required this.action,
    this.direction,
    this.readonly = false,
  }) : super(key: key);

  @override
  _GameLocationGridWidgetState createState() => _GameLocationGridWidgetState();
}

typedef DungeonGridMemberFunction = Widget Function(
    DungeonActionRecord record, String key);

class _GameLocationGridWidgetState extends State<GameLocationGridWidget> {
  double gridMemberWidth = 0;
  double gridMemberHeight = 0;

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

  @override
  void initState() {
    super.initState();
  }

  List<Widget> _generateGrid(BuildContext context) {
    var locationContents = getLocationContents(widget.locationData);

    int roomGridIdx = 0;
    List<Widget Function()> dunegonGridMemberFunctions = [
      // Top Row
      () => _directionWidget(context, widget.locationData, 'northwest'),
      () => _roomWidget(context, locationContents, roomGridIdx++),
      () => _directionWidget(context, widget.locationData, 'north'),
      () => _roomWidget(context, locationContents, roomGridIdx++),
      () => _directionWidget(context, widget.locationData, 'northeast'),
      // Second Row
      () => _roomWidget(context, locationContents, roomGridIdx++),
      () => _roomWidget(context, locationContents, roomGridIdx++),
      () => _directionWidget(context, widget.locationData, 'up'),
      () => _roomWidget(context, locationContents, roomGridIdx++),
      () => _roomWidget(context, locationContents, roomGridIdx++),
      // Third Row
      () => _directionWidget(context, widget.locationData, 'west'),
      () => _roomWidget(context, locationContents, roomGridIdx++),
      () => _roomWidget(context, locationContents, roomGridIdx++),
      () => _roomWidget(context, locationContents, roomGridIdx++),
      () => _directionWidget(context, widget.locationData, 'east'),
      // Fourth Row
      () => _roomWidget(context, locationContents, roomGridIdx++),
      () => _roomWidget(context, locationContents, roomGridIdx++),
      () => _directionWidget(context, widget.locationData, 'down'),
      () => _roomWidget(context, locationContents, roomGridIdx++),
      () => _roomWidget(context, locationContents, roomGridIdx++),
      // Bottom Row
      () => _directionWidget(context, widget.locationData, 'southwest'),
      () => _roomWidget(context, locationContents, roomGridIdx++),
      () => _directionWidget(context, widget.locationData, 'south'),
      () => _roomWidget(context, locationContents, roomGridIdx++),
      () => _directionWidget(context, widget.locationData, 'southeast'),
    ];

    List<Widget> gridWidgets = [];
    for (var gridMemberFunction in dunegonGridMemberFunctions) {
      gridWidgets.add(gridMemberFunction());
    }

    return gridWidgets;
  }

  // Direction widget
  Widget _directionWidget(
      BuildContext context, LocationData locationData, String direction) {
    if (!locationData.directions.contains(direction)) {
      return _emptyWidget('${directionLabelMap[direction]}');
    }

    return Container(
      margin: const EdgeInsets.all(2),
      child: ElevatedButton(
        onPressed: () {
          _selectTarget(context, direction);
        },
        style: gameButtonStyle,
        child: Text('${directionLabelMap[direction]}'),
      ),
    );
  }

  // Room widget
  Widget _roomWidget(BuildContext context,
      Map<int, LocationContent> locationContents, int idx) {
    if (locationContents[idx] == null) {
      return _emptyWidget('E$idx');
    }
    Widget returnWidget;
    var locationContent = locationContents[idx];
    switch (locationContent!.type) {
      case ContentType.character:
        {
          returnWidget = _characterWidget(context, locationContent.name);
        }
        break;
      case ContentType.monster:
        {
          returnWidget = _monsterWidget(context, locationContent.name);
        }
        break;
      case ContentType.object:
        {
          returnWidget = _objectWidget(context, locationContent.name);
        }
        break;
      default:
        {
          returnWidget = _emptyWidget('E$idx');
        }
    }
    return returnWidget;
  }

  // Character widget
  Widget _characterWidget(BuildContext context, String characterName) {
    return Container(
      margin: gameButtonMargin,
      child: ElevatedButton(
        onPressed: () {
          final log = getLogger('GameLocationGridWidget');
          log.info('Selecting character >$characterName<');
          _selectTarget(context, characterName);
        },
        style: gameButtonStyle,
        child: Text(characterName),
      ),
    );
  }

  // Monster widget
  Widget _monsterWidget(BuildContext context, String monsterName) {
    return Container(
      margin: gameButtonMargin,
      child: ElevatedButton(
        onPressed: () {
          final log = getLogger('GameLocationGridWidget');
          log.info('Selecting monster >$monsterName<');
          _selectTarget(context, monsterName);
        },
        style: gameButtonStyle,
        child: Text(monsterName),
      ),
    );
  }

  // Object widget
  Widget _objectWidget(BuildContext context, String objectName) {
    return Container(
      margin: gameButtonMargin,
      child: ElevatedButton(
        onPressed: () {
          final log = getLogger('GameLocationGridWidget');
          log.info('Selecting object >$objectName<');
          _selectTarget(context, objectName);
        },
        style: gameButtonStyle,
        child: Text(objectName),
      ),
    );
  }

  // Empty widget
  Widget _emptyWidget(String label) {
    return Container(
      width: gridMemberWidth,
      height: gridMemberHeight,
      alignment: Alignment.center,
      margin: gameButtonMargin,
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

  void _selectTarget(BuildContext context, String target) {
    final log = getLogger('GameLocationGridWidget');
    log.info('Submitting move action..');

    final dungeonCubit = BlocProvider.of<DungeonCubit>(context);
    if (dungeonCubit.dungeonRecord == null) {
      log.warning(
          'Dungeon cubit missing dungeon record, cannot initialise action');
      return;
    }

    final characterCubit = BlocProvider.of<CharacterCubit>(context);
    if (characterCubit.characterRecord == null) {
      log.warning(
          'Character cubit missing character record, cannot initialise action');
      return;
    }

    final dungeonCommandCubit = BlocProvider.of<DungeonCommandCubit>(context);
    if (dungeonCommandCubit.target == target) {
      log.info('++ Unselecting target $target');
      dungeonCommandCubit.unselectTarget();
      return;
    }

    log.info('++ Selecting target $target');
    dungeonCommandCubit.selectTarget(target);
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameLocationGridWidget');

    return LayoutBuilder(
      builder: (BuildContext context, BoxConstraints constraints) {
        log.info(
            'Building width ${constraints.maxWidth} height ${constraints.maxHeight}');

        // Set grid member dimensions
        gridMemberWidth = (constraints.maxWidth / 5) - 2;
        gridMemberHeight = (constraints.maxHeight / 5) - 2;
        if (gridMemberHeight > gridMemberWidth) {
          gridMemberHeight = gridMemberWidth;
        }
        if (gridMemberWidth > gridMemberHeight) {
          gridMemberWidth = gridMemberHeight;
        }

        log.info(
          '(B-**) Resulting button width $gridMemberWidth height $gridMemberHeight',
        );

        return IgnorePointer(
          ignoring: widget.readonly ? true : false,
          child: Container(
            margin: const EdgeInsets.fromLTRB(4, 4, 4, 4),
            decoration: BoxDecoration(
              color: widget.readonly ? Colors.black : const Color(0xFFDEDEDE),
              border: Border.all(
                color: const Color(0xFFDEDEDE),
              ),
              borderRadius: const BorderRadius.all(Radius.circular(5)),
            ),
            child: GridView.count(
              crossAxisCount: 5,
              children: _generateGrid(context),
            ),
          ),
        );
      },
    );
  }
}
