import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/repository/dungeon_action/dungeon_action_repository.dart';
import 'package:go_mud_client/widgets/game/game_dungeon_grid.dart';

class GameDungeonLookGridWidget extends StatefulWidget {
  final LocationData locationData;
  final String? action;
  final String? direction;

  const GameDungeonLookGridWidget({
    Key? key,
    required this.locationData,
    this.action,
    this.direction,
  }) : super(key: key);

  @override
  _GameDungeonLookGridWidgetState createState() => _GameDungeonLookGridWidgetState();
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

class _GameDungeonLookGridWidgetState extends State<GameDungeonLookGridWidget>
    with SingleTickerProviderStateMixin {
  // Animation controller
  late final AnimationController _controller = AnimationController(
    duration: const Duration(milliseconds: 250),
    vsync: this,
  );

  // Animation
  late final Animation<Offset> _offsetAnimation;

  double gridMemberWidth = 0;
  double gridMemberHeight = 0;

  double _opacity = 0.0;
  int _milliseconds = 1;

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
    final log = getLogger('GameDungeonLookGridWidget');

    Offset beginOffset = Offset.zero;
    Offset endOffset = Offset.zero;

    log.info('(initState) Target dungeon location direction ${widget.direction}');

    if (widget.direction != null) {
      beginOffset = slideInBeginOffset[widget.direction]!;
      endOffset = Offset.zero;
    }

    _controller.addStatusListener((status) {
      if (status == AnimationStatus.forward) {
        setState(() {
          _opacity = 1.0;
          _milliseconds = 500;
        });
      }
      if (status == AnimationStatus.completed) {
        Future.delayed(const Duration(milliseconds: 1500), () {
          setState(() {
            _opacity = 0.0;
            _milliseconds = 500;
          });
        });
      }
    });

    _offsetAnimation = Tween<Offset>(
      begin: beginOffset,
      end: endOffset,
    ).animate(CurvedAnimation(
      parent: _controller,
      curve: Curves.linear,
    ));

    WidgetsBinding.instance?.addPostFrameCallback((_) => _controller.forward());

    super.initState();
  }

  @override
  void dispose() {
    super.dispose();
    _controller.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameDungeonLookGridWidget');
    log.info('Building ${widget.key}');

    return AnimatedBuilder(
      animation: _controller,
      child: AnimatedOpacity(
        opacity: _opacity,
        duration: Duration(milliseconds: _milliseconds),
        child: GameDungeonGridWidget(
          locationData: widget.locationData,
          action: widget.action,
          readonly: true,
        ),
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
