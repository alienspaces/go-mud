import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/widgets/game/action/action.dart';
import 'package:go_mud_client/cubit/dungeon_command/dungeon_command_cubit.dart';

class GameActionCommandWidget extends StatefulWidget {
  const GameActionCommandWidget({Key? key}) : super(key: key);

  @override
  State<GameActionCommandWidget> createState() =>
      _GameActionCommandWidgetState();
}

// TODO: 9-implement-monster-actions: Timer with animation that submits a
// look action an action is not performed by the player.

class _GameActionCommandWidgetState extends State<GameActionCommandWidget>
    with SingleTickerProviderStateMixin {
  late Animation<double> animation;
  late AnimationController controller;

  @override
  void initState() {
    super.initState();

    controller =
        AnimationController(vsync: this, duration: const Duration(seconds: 8));

    animation = Tween<double>(begin: 0, end: 100).animate(controller);
    animation.addListener(() {
      final log = getLogger('Command', 'animationListener');
      final value = animation.value.toInt();
      if (value == 99) {
        log.warning("** Stopping starting animation");
        submitLookAction(context);
        _restartAnimation();
      }
      setState(() {});
    });

    animation.addStatusListener((status) {
      final log = getLogger('Command', 'animationStatusListener');
      log.warning("** Animation status $status");
      if (status == AnimationStatus.forward) {
        // submitLookAction(context);
      }
    });

    controller.forward();
    submitLookAction(context);
  }

  _restartAnimation() {
    controller.stop();
    controller.forward(from: 0);
  }

  @override
  void dispose() {
    controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final maxWidth = MediaQuery.of(context).size.width;
    final width = maxWidth - ((maxWidth * animation.value) / 100);

    return BlocConsumer<DungeonCommandCubit, DungeonCommandState>(
      listener: (BuildContext context, DungeonCommandState state) {
        final log = getLogger('Command', 'buildListener');
        if (state is DungeonCommandStatePreparing) {
          log.warning(
              "** command state preparing ${state.action} ${state.target}");
          _restartAnimation();
        }
      },
      builder: (BuildContext context, DungeonCommandState state) {
        // ignore: avoid_unnecessary_containers
        Widget commandWidget = Container(child: const Text(" "));

        if (state is DungeonCommandStatePreparing) {
          commandWidget = Container(
            color: Colors.brown[200],
            alignment: Alignment.center,
            child:
                Text('${state.action ?? ''} ${state.target ?? ''}'.trimRight()),
          );
        }

        return Stack(children: [
          commandWidget,
          Container(
            color: Colors.green[200]!.withOpacity(0.5),
            width: width,
          ),
        ]);
      },
    );
  }
}
