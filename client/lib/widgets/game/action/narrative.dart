import 'package:flutter/material.dart';
// import 'package:flutter/scheduler.dart';
import 'dart:async';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';

class GameActionNarrativeWidget extends StatefulWidget {
  const GameActionNarrativeWidget({Key? key}) : super(key: key);
  @override
  _GameActionNarrativeWidgetState createState() =>
      _GameActionNarrativeWidgetState();
}

class _GameActionNarrativeWidgetState extends State<GameActionNarrativeWidget> {
  List<Widget> lines = [];

  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameActionNarrativeWidget');
    log.info('Building..');

    return BlocConsumer<DungeonActionCubit, DungeonActionState>(
      listener: (BuildContext context, DungeonActionState state) {
        log.info('listener...');
        if (state is DungeonActionStatePlaying) {
          // Add narrative line
          var lineWidget = GameActionNarrativeLineWidget(
            line: state.current.narrative,
          );
          lines.add(lineWidget);
        }
      },
      builder: (BuildContext context, DungeonActionState state) {
        // ignore: avoid_unnecessary_containers
        return IgnorePointer(
          ignoring: true,
          child: Stack(
            children: lines,
          ),
        );
      },
    );
  }
}

class GameActionNarrativeLineWidget extends StatefulWidget {
  final String line;

  const GameActionNarrativeLineWidget({Key? key, required this.line})
      : super(key: key);
  @override
  _GameActionNarrativeLineWidgetState createState() =>
      _GameActionNarrativeLineWidgetState();
}

class _GameActionNarrativeLineWidgetState
    extends State<GameActionNarrativeLineWidget> {
  double opacity = 1.0;
  double bottom = 0.0;
  late Timer animationTimer;
  @override
  initState() {
    final log = getLogger('GameActionNarrativeLineWidget');
    super.initState();
    log.info('### Initialised state');

    animationTimer = Timer(const Duration(milliseconds: 200), () {
      double newOpacity = 0.0;
      double newBottom = 1000;
      setState(() {
        opacity = newOpacity;
        bottom = newBottom;
      });
      log.info('^^^ Updated opacity "$newOpacity" bottom "$newBottom"');
    });
  }

  @override
  void dispose() {
    animationTimer.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameActionNarrativeLineWidget');
    log.info(
        '^^^ Building with line "${widget.line}" opacity "$opacity" bottom "$bottom"');

    return AnimatedPositioned(
      bottom: bottom,
      duration: const Duration(milliseconds: 2500),
      child: Container(
        alignment: Alignment.center,
        child: AnimatedOpacity(
          opacity: opacity,
          duration: const Duration(milliseconds: 2500),
          child: Text(': ${widget.line}'.trimRight()),
        ),
      ),
    );
  }
}
