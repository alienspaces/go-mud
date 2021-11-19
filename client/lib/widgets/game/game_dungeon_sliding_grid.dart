import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/repository/dungeon_action/dungeon_action_repository.dart';
import 'package:go_mud_client/widgets/game/game_dungeon_grid.dart';

enum Slide { slideIn, slideOut, slideNone }

class GameDungeonSlidingGridWidget extends StatefulWidget {
  final Slide slide;
  final DungeonActionRecord dungeonActionRecord;
  final String action;
  final String? direction;

  const GameDungeonSlidingGridWidget({
    Key? key,
    required this.slide,
    required this.dungeonActionRecord,
    required this.action,
    this.direction,
  }) : super(key: key);

  @override
  _GameDungeonSlidingGridWidgetState createState() => _GameDungeonSlidingGridWidgetState();
}

Map<String, Offset> slideInBeginOffset = {
  'north': const Offset(0, -1),
  'northeast': const Offset(1, -1),
  'east': const Offset(1, 0),
  'southeast': const Offset(1, 1),
  'south': const Offset(0, 1),
  'southwest': const Offset(-1, 1),
  'west': const Offset(-1, 0),
  'northwest': const Offset(-1, -1),
  'up': const Offset(-.1, -1),
  'down': const Offset(.1, 1),
};

Map<String, Offset> slideOutEndOffset = {
  'north': const Offset(0, 1),
  'northeast': const Offset(-1, 1),
  'east': const Offset(-1, 0), //
  'southeast': const Offset(-1, -1),
  'south': const Offset(0, -1),
  'southwest': const Offset(1, -1),
  'west': const Offset(1, 0),
  'northwest': const Offset(1, 1),
  'up': const Offset(.1, 1),
  'down': const Offset(-.1, -1),
};

class _GameDungeonSlidingGridWidgetState extends State<GameDungeonSlidingGridWidget>
    with SingleTickerProviderStateMixin {
  // Animation controller
  late final AnimationController _controller;
  // Animation
  late final Animation<Offset> _offsetAnimation;

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
    final log = getLogger('GameDungeonSlidingGridWidget');

    // Animation controller
    _controller = AnimationController(
      duration: const Duration(milliseconds: 500),
      vsync: this,
    );

    Offset beginOffset = Offset.zero;
    Offset endOffset = Offset.zero;

    String command = widget.dungeonActionRecord.action.command;
    log.info('(initState) Target dungeon action id ${widget.dungeonActionRecord.action.id}');
    log.info('(initState) Target dungeon location command $command');
    log.info('(initState) Target dungeon location direction ${widget.direction}');
    log.info('(initState) Target dungeon location slide ${widget.slide}');

    if (command == 'move' && widget.direction != null) {
      if (widget.slide == Slide.slideIn) {
        beginOffset = slideInBeginOffset[widget.direction]!;
        endOffset = Offset.zero;
      } else if (widget.slide == Slide.slideOut) {
        beginOffset = Offset.zero;
        endOffset = slideOutEndOffset[widget.direction]!;
      }
    }

    _offsetAnimation = Tween<Offset>(
      begin: beginOffset,
      end: endOffset,
    ).animate(CurvedAnimation(
      parent: _controller,
      curve: Curves.linear,
    ));

    if (widget.slide != Slide.slideNone) {
      WidgetsBinding.instance?.addPostFrameCallback((_) => _controller.forward());
    }

    super.initState();
  }

  @override
  void dispose() {
    super.dispose();
    _controller.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameDungeonSlidingGridWidget');
    log.info('Building..');

    return AnimatedBuilder(
      animation: _controller,
      child: GameDungeonGridWidget(
        dungeonActionRecord: widget.dungeonActionRecord,
        action: widget.action,
      ),
      builder: (BuildContext context, Widget? child) {
        return SlideTransition(
          position: _offsetAnimation,
          child: child,
        );
      },
    );
  }
}
