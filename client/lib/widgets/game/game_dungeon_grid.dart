import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/location.dart';
import 'package:go_mud_client/cubit/dungeon/dungeon_cubit.dart';
// import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_command/dungeon_command_cubit.dart';
import 'package:go_mud_client/cubit/character/character_cubit.dart';
import 'package:go_mud_client/repository/dungeon_action/dungeon_action_repository.dart';

enum DungeonGridScroll { scrollIn, scrollOut, scrollNone }

class GameDungeonGridWidget extends StatefulWidget {
  final DungeonGridScroll scroll;
  final DungeonActionRecord dungeonActionRecord;
  final String? direction;

  const GameDungeonGridWidget(
      {Key? key, required this.scroll, required this.dungeonActionRecord, this.direction})
      : super(key: key);

  @override
  _GameDungeonGridWidgetState createState() => _GameDungeonGridWidgetState();
}

typedef DungeonGridMemberFunction = Widget Function(DungeonActionRecord record, String key);

Map<String, Offset> scrollInBeginOffset = {
  'north': const Offset(0, -1),
  'northeast': const Offset(1, -1),
  'east': const Offset(1, 0),
  'southeast': const Offset(1, 1),
  'south': const Offset(0, 1),
  'southwest': const Offset(-1, 1),
  'west': const Offset(-1, 0),
  'northwest': const Offset(-1, -1),
};

Map<String, Offset> scrollOutEndOffset = {
  'north': const Offset(0, 1),
  'northeast': const Offset(-1, 1),
  'east': const Offset(-1, 0), //
  'southeast': const Offset(-1, -1),
  'south': const Offset(0, -1),
  'southwest': const Offset(1, -1),
  'west': const Offset(1, 0),
  'northwest': const Offset(1, 1),
};

class _GameDungeonGridWidgetState extends State<GameDungeonGridWidget>
    with SingleTickerProviderStateMixin {
  late final AnimationController _controller;
  // Animation controller
  late final Animation<Offset> _offsetAnimation;
  @override
  void initState() {
    final log = getLogger('DungeonActionCubit');

    // Animation controller
    _controller = AnimationController(
      duration: const Duration(seconds: 1),
      vsync: this,
    );

    Offset beginOffset = Offset.zero;
    Offset endOffset = Offset.zero;

    String command = widget.dungeonActionRecord.action.command;
    log.info('(initState) Target dungeon location command $command');
    log.info('(initState) Target dungeon location direction ${widget.direction}');

    if (command == 'move' && widget.direction != null) {
      if (widget.scroll == DungeonGridScroll.scrollIn) {
        beginOffset = scrollInBeginOffset[widget.direction]!;
        endOffset = Offset.zero;
      } else if (widget.scroll == DungeonGridScroll.scrollOut) {
        beginOffset = Offset.zero;
        endOffset = scrollOutEndOffset[widget.direction]!;
      }
    }

    _offsetAnimation = Tween<Offset>(
      begin: beginOffset,
      end: endOffset,
    ).animate(CurvedAnimation(
      parent: _controller,
      curve: Curves.linear,
    ));

    super.initState();
  }

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

  List<Widget> _generateGrid(BuildContext context) {
    var dungeonActionRecord = widget.dungeonActionRecord;

    var locationContents = getLocationContents(dungeonActionRecord);

    int roomGridIdx = 0;
    List<Widget Function()> dunegonGridMemberFunctions = [
      // Top Row
      () => _directionWidget(context, dungeonActionRecord, 'northwest'),
      () => _roomWidget(context, locationContents, roomGridIdx++),
      () => _directionWidget(context, dungeonActionRecord, 'north'),
      () => _roomWidget(context, locationContents, roomGridIdx++),
      () => _directionWidget(context, dungeonActionRecord, 'northeast'),
      // Second Row
      () => _roomWidget(context, locationContents, roomGridIdx++),
      () => _roomWidget(context, locationContents, roomGridIdx++),
      () => _directionWidget(context, dungeonActionRecord, 'up'),
      () => _roomWidget(context, locationContents, roomGridIdx++),
      () => _roomWidget(context, locationContents, roomGridIdx++),
      // Third Row
      () => _directionWidget(context, dungeonActionRecord, 'west'),
      () => _roomWidget(context, locationContents, roomGridIdx++),
      () => _roomWidget(context, locationContents, roomGridIdx++),
      () => _roomWidget(context, locationContents, roomGridIdx++),
      () => _directionWidget(context, dungeonActionRecord, 'east'),
      // Fourth Row
      () => _roomWidget(context, locationContents, roomGridIdx++),
      () => _roomWidget(context, locationContents, roomGridIdx++),
      () => _directionWidget(context, dungeonActionRecord, 'down'),
      () => _roomWidget(context, locationContents, roomGridIdx++),
      () => _roomWidget(context, locationContents, roomGridIdx++),
      // Bottom Row
      () => _directionWidget(context, dungeonActionRecord, 'southwest'),
      () => _roomWidget(context, locationContents, roomGridIdx++),
      () => _directionWidget(context, dungeonActionRecord, 'south'),
      () => _roomWidget(context, locationContents, roomGridIdx++),
      () => _directionWidget(context, dungeonActionRecord, 'southeast'),
    ];

    List<Widget> gridWidgets = [];
    for (var gridMemberFunction in dunegonGridMemberFunctions) {
      gridWidgets.add(gridMemberFunction());
    }

    return gridWidgets;
  }

  // Direction widget
  Widget _directionWidget(BuildContext context, DungeonActionRecord record, String direction) {
    if (!record.location.directions.contains(direction)) {
      return _emptyWidget('${directionLabelMap[direction]}');
    }

    return Container(
      margin: const EdgeInsets.all(2),
      child: ElevatedButton(
        onPressed: () {
          _selectTarget(context, direction);
        },
        child: Text('${directionLabelMap[direction]}'),
      ),
    );
  }

  // Room widget
  Widget _roomWidget(BuildContext context, Map<int, LocationContent> locationContents, int idx) {
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
      margin: const EdgeInsets.all(2),
      child: ElevatedButton(
        onPressed: () {
          final log = getLogger('GameDungeonGridWidget');
          log.info('Selecting character >$characterName<');
          _selectTarget(context, characterName);
        },
        style: ElevatedButton.styleFrom(
          primary: Colors.green,
        ),
        child: Text(characterName),
      ),
    );
  }

  // Monster widget
  Widget _monsterWidget(BuildContext context, String monsterName) {
    return Container(
      margin: const EdgeInsets.all(2),
      child: ElevatedButton(
        onPressed: () {
          final log = getLogger('GameDungeonGridWidget');
          log.info('Selecting monster >$monsterName<');
          _selectTarget(context, monsterName);
        },
        style: ElevatedButton.styleFrom(
          primary: Colors.orange,
        ),
        child: Text(monsterName),
      ),
    );
  }

  // Object widget
  Widget _objectWidget(BuildContext context, String objectName) {
    return Container(
      margin: const EdgeInsets.all(2),
      child: ElevatedButton(
        onPressed: () {
          final log = getLogger('GameDungeonGridWidget');
          log.info('Selecting object >$objectName<');
          _selectTarget(context, objectName);
        },
        style: ElevatedButton.styleFrom(
          primary: Colors.brown,
        ),
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

  void _selectTarget(BuildContext context, String target) {
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
  void dispose() {
    super.dispose();
    _controller.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameDungeonGridWidget');
    log.info('Building..');

    if (widget.scroll != DungeonGridScroll.scrollNone) {
      _controller.forward();
    }

    double gridWidth = gridMemberWidth * 5;
    double gridHeight = gridMemberHeight * 5;

    return SlideTransition(
      position: _offsetAnimation,
      child: Container(
        decoration: BoxDecoration(
          color: const Color(0xFFDEDEDE),
          border: Border.all(
            color: const Color(0xFFDEDEDE),
          ),
          borderRadius: const BorderRadius.all(Radius.circular(5)),
        ),
        padding: const EdgeInsets.all(1),
        margin: const EdgeInsets.all(5),
        width: gridWidth,
        height: gridHeight,
        child: GridView.count(
          crossAxisCount: 5,
          children: _generateGrid(context),
        ),
      ),
    );
  }
}
