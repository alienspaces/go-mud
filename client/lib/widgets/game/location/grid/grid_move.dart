import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/repository/dungeon_action/dungeon_action_repository.dart';
import 'package:go_mud_client/widgets/game/location/grid/grid.dart';

enum Slide { slideIn, slideOut, slideNone }

class GameLocationGridMoveWidget extends StatefulWidget {
  final Slide slide;
  final LocationData locationData;
  final String? action;
  final String? direction;

  const GameLocationGridMoveWidget({
    Key? key,
    required this.slide,
    required this.locationData,
    this.action,
    this.direction,
  }) : super(key: key);

  @override
  _GameLocationGridMoveWidgetState createState() =>
      _GameLocationGridMoveWidgetState();
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

class _GameLocationGridMoveWidgetState extends State<GameLocationGridMoveWidget>
    with SingleTickerProviderStateMixin {
  // Animation controller
  late final AnimationController _controller = AnimationController(
    duration: const Duration(milliseconds: 500),
    vsync: this,
  );

  // Animation
  late final Animation<Offset> _offsetAnimation;

  @override
  void initState() {
    final log = getLogger('GameLocationGridMoveWidget');

    Offset beginOffset = Offset.zero;
    Offset endOffset = Offset.zero;

    log.info(
        '(initState) Target dungeon location direction ${widget.direction}');
    log.info('(initState) Target dungeon location slide ${widget.slide}');

    if (widget.direction != null) {
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
      WidgetsBinding.instance
          ?.addPostFrameCallback((_) => _controller.forward());
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
    final log = getLogger('GameLocationGridMoveWidget');
    log.info('Building ${widget.key} - ${widget.slide}..');

    return AnimatedBuilder(
      animation: _controller,
      child: GameLocationGridWidget(
        locationData: widget.locationData,
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
