import 'package:flutter/material.dart';
// import 'package:flutter/scheduler.dart';
import 'dart:async';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';

class NarrativeLineFormat {
  final Color colour;
  const NarrativeLineFormat({required this.colour});
}

const NarrativeLineFormat characterNarrativeLineFormat = NarrativeLineFormat(
  colour: Colors.brown,
);

const NarrativeLineFormat otherNarrativeLineFormat = NarrativeLineFormat(
  colour: Colors.orange,
);

const NarrativeLineFormat errorNarrativeLineFormat = NarrativeLineFormat(
  colour: Colors.red,
);

class GameActionNarrativeWidget extends StatefulWidget {
  const GameActionNarrativeWidget({Key? key}) : super(key: key);
  @override
  State<GameActionNarrativeWidget> createState() =>
      _GameActionNarrativeWidgetState();
}

class _GameActionNarrativeWidgetState extends State<GameActionNarrativeWidget> {
  List<Widget> lines = [];

  @override
  void dispose() {
    final log = getLogger('GameActionNarrativeWidget', 'dispose');
    if (!mounted) {
      log.info('### Not mounted, not disposing..');
      return;
    }
    log.info('### Disposing...');
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameActionNarrativeWidget', 'build');
    log.info('Building..');
    double width = MediaQuery.of(context).size.width;

    return BlocConsumer<DungeonActionCubit, DungeonActionState>(
      listener: (BuildContext context, DungeonActionState state) {
        log.fine('listener... width $width');
        if (state is DungeonActionStatePlaying) {
          // Add narrative line
          var lineWidget = GameActionNarrativeLineWidget(
            key: UniqueKey(),
            line: state.currentActionRec.actionNarrative,
            format: characterNarrativeLineFormat,
            width: width,
          );
          setState(() {
            lines.add(lineWidget);
          });
        } else if (state is DungeonActionStatePlayingOther) {
          // Add narrative line
          var lineWidget = GameActionNarrativeLineWidget(
            key: UniqueKey(),
            line: state.actionRec.actionNarrative,
            format: otherNarrativeLineFormat,
            width: width,
          );
          setState(() {
            lines.add(lineWidget);
          });
        } else if (state is DungeonActionStateError) {
          var lineWidget = GameActionNarrativeLineWidget(
            key: UniqueKey(),
            line: state.message,
            format: errorNarrativeLineFormat,
            width: width,
          );
          setState(() {
            lines.add(lineWidget);
          });
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
  final NarrativeLineFormat format;
  final double width;

  const GameActionNarrativeLineWidget({
    Key? key,
    required this.line,
    required this.format,
    required this.width,
  }) : super(key: key);
  @override
  State<GameActionNarrativeLineWidget> createState() =>
      _GameActionNarrativeLineWidgetState();
}

class _GameActionNarrativeLineWidgetState
    extends State<GameActionNarrativeLineWidget> {
  double opacity = 1.0;
  double bottom = 0.0;
  late Timer animationTimer;
  @override
  initState() {
    final log = getLogger('GameActionNarrativeLineWidget', 'initState');

    super.initState();

    if (!mounted) {
      log.info('### Not mounted, not initialising..');
      return;
    }

    log.info('### Initialising..');

    animationTimer = Timer(const Duration(milliseconds: 200), () {
      double newOpacity = 0.0;
      double newBottom = 1500;
      setState(() {
        opacity = newOpacity;
        bottom = newBottom;
      });
      log.info('^^^ Updated opacity "$newOpacity" bottom "$newBottom"');
    });
  }

  @override
  void dispose() {
    final log = getLogger('GameActionNarrativeLineWidget', 'dispose');

    if (!mounted) {
      log.info('### Not mounted, not disposing..');
      return;
    }

    log.info('### Disposing..');
    animationTimer.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameActionNarrativeLineWidget', 'build');

    if (!mounted) {
      log.info('^^^ Not mounted, not building with line "${widget.line}"');
      return Container();
    }

    log.info(
        '^^^ Building with line "${widget.line}" opacity "$opacity" bottom "$bottom"');

    return AnimatedPositioned(
      bottom: bottom,
      duration: const Duration(milliseconds: 3000),
      child: Container(
        width: widget.width,
        alignment: Alignment.center,
        child: AnimatedOpacity(
          opacity: opacity,
          duration: const Duration(milliseconds: 1500),
          child: Container(
            margin: const EdgeInsets.all(3),
            alignment: Alignment.center,
            child: Text(
              ': ${widget.line}'.trimRight(),
              textAlign: TextAlign.center,
              style: Theme.of(context).textTheme.titleMedium!.copyWith(
                    fontSize: 30,
                    color: widget.format.colour,
                  ),
            ),
          ),
        ),
      ),
    );
  }
}
