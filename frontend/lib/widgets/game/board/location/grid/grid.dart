import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/location.dart';
import 'package:go_mud_client/style.dart';
import 'package:go_mud_client/repository/dungeon_action/dungeon_action_repository.dart';
import 'package:go_mud_client/widgets/game/button/character_button.dart';
import 'package:go_mud_client/widgets/game/button/monster_button.dart';
import 'package:go_mud_client/widgets/game/button/object_button.dart';
import 'package:go_mud_client/cubit/target.dart';

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
  State<GameLocationGridWidget> createState() => _GameLocationGridWidgetState();
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
    if (!locationData.locationDirections.contains(direction)) {
      return _emptyWidget('${directionLabelMap[direction]}');
    }

    return Container(
      margin: const EdgeInsets.all(2),
      child: ElevatedButton(
        onPressed: () {
          _selectTarget(context, direction);
        },
        style: gameBoardButtonStyle,
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
    if (locationContent == null) {
      return _emptyWidget('E$idx');
    }

    switch (locationContent.type) {
      case ContentType.character:
        {
          returnWidget = _characterWidget(context, locationContent);
        }
        break;
      case ContentType.monster:
        {
          returnWidget = _monsterWidget(context, locationContent);
        }
        break;
      case ContentType.object:
        {
          returnWidget = _objectWidget(context, locationContent);
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
  Widget _characterWidget(BuildContext context, LocationContent character) {
    return CharacterButtonWidget(
      name: character.name,
      health: character.health ?? 0,
      currentHealth: character.currentHealth ?? 0,
      fatigue: character.fatigue ?? 0,
      currentFatigue: character.currentFatigue ?? 0,
    );
  }

  // Monster widget
  Widget _monsterWidget(BuildContext context, LocationContent monster) {
    return MonsterButtonWidget(
      name: monster.name,
      health: monster.health ?? 0,
      currentHealth: monster.currentHealth ?? 0,
      fatigue: monster.fatigue ?? 0,
      currentFatigue: monster.currentFatigue ?? 0,
    );
  }

  // Object widget
  Widget _objectWidget(BuildContext context, LocationContent object) {
    return ObjectButtonWidget(name: object.name);
  }

  // Empty widget
  Widget _emptyWidget(String label) {
    return Container(
      width: gridMemberWidth,
      height: gridMemberHeight,
      alignment: Alignment.center,
      margin: gameBoardButtonMargin,
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
    selectTarget(context, target);
  }

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (BuildContext context, BoxConstraints constraints) {
        // Set grid member dimensions
        gridMemberWidth = (constraints.maxWidth / 5) - 2;
        gridMemberHeight = (constraints.maxHeight / 5) - 2;
        if (gridMemberHeight > gridMemberWidth) {
          gridMemberHeight = gridMemberWidth;
        }
        if (gridMemberWidth > gridMemberHeight) {
          gridMemberWidth = gridMemberHeight;
        }

        return IgnorePointer(
          ignoring: widget.readonly ? true : false,
          child: Container(
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
