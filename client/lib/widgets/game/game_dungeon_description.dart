import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/repository/dungeon_action/dungeon_action_repository.dart';

enum DescriptionOpacity { fadeIn, fadeOut }

class GameDungeonDescriptionWidget extends StatefulWidget {
  final DescriptionOpacity fade;
  final DungeonActionRecord dungeonActionRecord;

  const GameDungeonDescriptionWidget(
      {Key? key, required this.fade, required this.dungeonActionRecord})
      : super(key: key);

  @override
  _GameDungeonDescriptionWidgetState createState() => _GameDungeonDescriptionWidgetState();
}

typedef DungeonDescriptionMemberFunction = Widget Function(DungeonActionRecord record, String key);

class _GameDungeonDescriptionWidgetState extends State<GameDungeonDescriptionWidget>
    with TickerProviderStateMixin {
  late final AnimationController _controller = AnimationController(
    duration: const Duration(milliseconds: 500),
    vsync: this,
    lowerBound: 0.0,
    upperBound: 1.0,
  );
  late final Animation<double> _animation = CurvedAnimation(
    parent: _controller,
    curve: Curves.easeIn,
  );

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameDungeonDescriptionWidget');
    log.info('Building.. ${widget.fade}');

    return BlocConsumer<DungeonActionCubit, DungeonActionState>(
      listener: (BuildContext context, DungeonActionState state) {
        log.info('listener...');
      },
      builder: (BuildContext context, DungeonActionState state) {
        widget.fade == DescriptionOpacity.fadeIn
            ? _controller.forward(from: 0.0)
            : _controller.reverse(from: 1.0);

        // ignore: avoid_unnecessary_containers
        return FadeTransition(
          opacity: _animation,
          child: Column(
            children: [
              Container(
                alignment: Alignment.center,
                margin: const EdgeInsets.fromLTRB(5, 10, 5, 5),
                child: Text(widget.dungeonActionRecord.location.name,
                    style: Theme.of(context).textTheme.headline5),
              ),
              Container(
                alignment: Alignment.center,
                margin: const EdgeInsets.fromLTRB(5, 5, 5, 10),
                child: Text(widget.dungeonActionRecord.location.description),
              ),
            ],
          ),
        );
      },
    );
  }
}
